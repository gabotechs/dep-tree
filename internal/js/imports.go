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
		Imports: orderedmap.NewOrderedMap[string, []string](),
		Errors:  make([]error, 0),
	}

	for _, stmt := range jsFile.Statements {
		var importPath string
		var names []string
		switch {
		case stmt == nil:
			continue
		case stmt.StaticImport != nil:
			importPath = stmt.StaticImport.Path
			names = gatherNamesFromStaticImport(stmt.StaticImport)
		case stmt.DynamicImport != nil:
			importPath = stmt.DynamicImport.Path
			names = []string{"*"}
		default:
			continue
		}
		var resolvedPath string
		ctx, resolvedPath, err = l.ResolvePath(ctx, importPath, path.Dir(filePath))
		if err != nil {
			result.Errors = append(result.Errors, err)
		} else if resolvedPath != "" {
			result.Imports.Set(resolvedPath, names)
		}
	}
	return ctx, result, nil
}

func gatherNamesFromStaticImport(si *grammar.StaticImport) []string {
	names := make([]string, 0)

	if imported := si.Imported; imported != nil {
		if imported.Default {
			names = append(names, "default")
		}
		if selection := imported.SelectionImport; selection != nil {
			if selection.AllImport != nil {
				names = append(names, "*")
			}
			if selection.Deconstruction != nil {
				names = append(names, selection.Deconstruction.Names...)
			}
		}
	} else {
		names = append(names, "*")
	}
	return names
}
