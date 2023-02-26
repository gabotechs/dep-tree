package js

import (
	"path"

	"dep-tree/internal/js/js_grammar"
	"dep-tree/internal/language"
)

func (l *Language) ParseImports(file *js_grammar.File) (*language.ImportsResult, error) {
	imports := make([]language.ImportEntry, 0)
	var errors []error

	for _, stmt := range file.Statements {
		importPath := ""
		entry := language.ImportEntry{}

		switch {
		case stmt == nil:
			continue
		case stmt.StaticImport != nil:
			importPath = stmt.StaticImport.Path
			if imported := stmt.StaticImport.Imported; imported != nil {
				if imported.Default {
					entry.Names = append(entry.Names, "default")
				}
				if selection := imported.SelectionImport; selection != nil {
					if selection.AllImport != nil {
						entry.All = true
					}
					if selection.Deconstruction != nil {
						entry.Names = append(entry.Names, selection.Deconstruction.Names...)
					}
				}
			} else {
				entry.All = true
			}
		case stmt.DynamicImport != nil:
			importPath = stmt.DynamicImport.Path
			entry.All = true
		default:
			continue
		}
		var err error
		entry.Id, err = l.ResolvePath(importPath, path.Dir(file.Path))
		if err != nil {
			errors = append(errors, err)
		} else if entry.Id != "" {
			imports = append(imports, entry)
		}
	}
	return &language.ImportsResult{
		Imports: imports,
		Errors:  errors,
	}, nil
}
