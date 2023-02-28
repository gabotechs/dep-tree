package rust

import (
	"dep-tree/internal/language"
	"dep-tree/internal/rust/rust_grammar"
)

func (l *Language) ParseExports(file *rust_grammar.File) (*language.ExportsResult, error) {
	exports := make([]language.ExportEntry, 0)
	var errors []error

	for _, stmt := range file.Statements {
		switch {
		case stmt.Use != nil && stmt.Use.Pub:
			names := make([]language.ExportName, len(stmt.Use.Names))
			for i, name := range stmt.Use.Names {
				names[i] = language.ExportName{Original: name.Original, Alias: name.Alias}
			}

			id, err := l.resolve(stmt.Use.PathSlices, file.Path)
			if err != nil {
				errors = append(errors, err)
				continue
			} else if id == "" {
				continue
			}
			exports = append(exports, language.ExportEntry{
				All:   stmt.Use.All,
				Names: names,
				Id:    id,
			})
		case stmt.Pub != nil:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{{Original: stmt.Pub.Name}},
				Id:    file.Path,
			})
		case stmt.Mod != nil && stmt.Mod.Pub:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{{Original: stmt.Mod.Name}},
				Id:    file.Path,
			})
		}
	}

	return &language.ExportsResult{
		Exports: exports,
		Errors:  errors,
	}, nil
}
