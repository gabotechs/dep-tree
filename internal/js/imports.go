package js

import (
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/js/js_grammar"
	"github.com/gabotechs/dep-tree/internal/language"
)

func (l *Language) ParseImports(file *language.FileInfo) (*language.ImportsResult, error) {
	imports := make([]language.ImportEntry, 0)
	var errors []error

	content := file.Content.(*js_grammar.File)
	for _, stmt := range content.Statements {
		importPath := ""
		entry := language.ImportEntry{}

		switch {
		case stmt == nil:
			continue
		case stmt.StaticImport != nil:
			importPath = stmt.StaticImport.Path
			if imported := stmt.StaticImport.Imported; imported != nil {
				if imported.Default {
					entry.Symbols = append(entry.Symbols, "default")
				}
				if selection := imported.SelectionImport; selection != nil {
					if selection.AllImport != nil {
						entry.All = true
					}
					if selection.Deconstruction != nil {
						entry.Symbols = append(entry.Symbols, selection.Deconstruction.Names...)
					}
				}
			} else {
				entry.All = true
			}
		case stmt.DynamicImport != nil:
			importPath = stmt.DynamicImport.Path
			entry.All = true
		case stmt.Require != nil:
			importPath = stmt.Require.Path
			entry.All = stmt.Require.Alias != ""
			entry.Symbols = stmt.Require.Names
		default:
			continue
		}
		var err error
		entry.AbsPath, err = l.ResolvePath(importPath, filepath.Dir(file.AbsPath))
		if err != nil {
			errors = append(errors, err)
		} else if entry.AbsPath != "" {
			imports = append(imports, entry)
		}
	}
	return &language.ImportsResult{
		Imports: imports,
		Errors:  errors,
	}, nil
}
