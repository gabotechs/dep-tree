package js

import (
	"context"
	"path"

	"github.com/elliotchance/orderedmap/v2"

	"dep-tree/internal/js/grammar"
	"dep-tree/internal/language"
)

type ImportsCacheKey string

func (l *Language) ParseImports(
	ctx context.Context,
	filePath string,
) (context.Context, *language.ImportsResult, error) {
	cacheKey := ImportsCacheKey(filePath)
	if cached, ok := ctx.Value(cacheKey).(*language.ImportsResult); ok {
		return ctx, cached, nil
	} else {
		ctx, result, err := l.uncachedParseImports(ctx, filePath)
		if err != nil {
			return ctx, nil, err
		}
		ctx = context.WithValue(ctx, cacheKey, result)
		return ctx, result, err
	}
}

func (l *Language) uncachedParseImports(
	ctx context.Context,
	filePath string,
) (context.Context, *language.ImportsResult, error) {
	ctx, jsFile, err := grammar.Parse(ctx, filePath)
	if err != nil {
		return ctx, nil, err
	}

	result := &language.ImportsResult{
		Imports: orderedmap.NewOrderedMap[string, language.ImportEntry](),
		Errors:  make([]error, 0),
	}

	for _, stmt := range jsFile.Statements {
		var importPath string

		entry := language.ImportEntry{}
		switch {
		case stmt == nil:
			continue
		case stmt.StaticImport != nil:
			importPath = stmt.StaticImport.Path
			entry = gatherNamesFromStaticImport(stmt.StaticImport)
		case stmt.DynamicImport != nil:
			importPath = stmt.DynamicImport.Path
			entry.All = true
		default:
			continue
		}
		var resolvedPath string
		ctx, resolvedPath, err = l.ResolvePath(ctx, importPath, path.Dir(filePath))
		if err != nil {
			result.Errors = append(result.Errors, err)
		} else if resolvedPath != "" {
			result.Imports.Set(resolvedPath, entry)
		}
	}
	return ctx, result, nil
}

func gatherNamesFromStaticImport(si *grammar.StaticImport) language.ImportEntry {
	entry := language.ImportEntry{}

	if imported := si.Imported; imported != nil {
		if imported.Default {
			entry.Names = append(entry.Names, "default")
		}
		if selection := imported.SelectionImport; selection != nil {
			if selection.AllImport != nil {
				return language.ImportEntry{All: true}
			}
			if selection.Deconstruction != nil {
				entry.Names = append(entry.Names, selection.Deconstruction.Names...)
			}
		}
	} else {
		return language.ImportEntry{All: true}
	}
	return entry
}
