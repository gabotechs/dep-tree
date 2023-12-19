package rust

import (
	"fmt"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/rust/rust_grammar"
)

func (l *Language) ParseExports(file *rust_grammar.File) (*language.ExportsEntries, error) {
	exports := make([]language.ExportEntry, 0)
	var errors []error

	for _, stmt := range file.Statements {
		switch {
		case stmt.Use != nil && stmt.Use.Pub:
			for _, use := range stmt.Use.Flatten() {
				path, err := l.resolve(use.PathSlices, file.Path)
				if err != nil {
					errors = append(errors, fmt.Errorf("error resolving use statement for name %s: %w", use.Name.Original, err))
					continue
				} else if path == "" {
					continue
				}

				if use.All {
					exports = append(exports, language.ExportEntry{
						All:  use.All,
						Path: path,
					})
				} else {
					exports = append(exports, language.ExportEntry{
						Names: []language.ExportName{{Original: use.Name.Original, Alias: use.Name.Alias}},
						Path:  path,
					})
				}
			}
		case stmt.Pub != nil:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{{Original: stmt.Pub.Name}},
				Path:  file.Path,
			})
		case stmt.Mod != nil && stmt.Mod.Pub:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{{Original: stmt.Mod.Name}},
				Path:  file.Path,
			})
		}
	}

	return &language.ExportsEntries{
		Exports: exports,
		Errors:  errors,
	}, nil
}
