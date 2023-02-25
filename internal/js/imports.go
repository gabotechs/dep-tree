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
	file *grammar.File,
) (context.Context, *language.ImportsResult, error) {
	result := &language.ImportsResult{
		Imports: orderedmap.NewOrderedMap[string, language.ImportEntry](),
		Errors:  make([]error, 0),
	}

	for _, stmt := range file.Statements {
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
		var err error
		ctx, resolvedPath, err = l.ResolvePath(ctx, importPath, path.Dir(file.Path))
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
