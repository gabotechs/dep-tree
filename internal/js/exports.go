package js

import (
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/js/js_grammar"
	"github.com/gabotechs/dep-tree/internal/language"
)

type ExportsCacheKey string

//nolint:gocyclo
func (l *Language) ParseExports(file *js_grammar.File) (*language.ExportsEntries, error) {
	exports := make([]language.ExportEntry, 0)
	var errors []error

	for _, stmt := range file.Statements {
		switch {
		case stmt == nil:
			// Is this even possible?
		case stmt.DeclarationExport != nil:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{
					{
						Original: stmt.DeclarationExport.Name,
					},
				},
				Path: file.Path,
			})
		case stmt.ListExport != nil:
			if stmt.ListExport.ExportDeconstruction != nil {
				for _, name := range stmt.ListExport.ExportDeconstruction.Names {
					exports = append(exports, language.ExportEntry{
						Names: []language.ExportName{
							{
								Original: name.Original,
								Alias:    name.Alias,
							},
						},
						Path: file.Path,
					})
				}
			}
		case stmt.DefaultExport != nil:
			if stmt.DefaultExport.Default {
				exports = append(exports, language.ExportEntry{
					Names: []language.ExportName{
						{
							Original: "default",
						},
					},
					Path: file.Path,
				})
			}
		case stmt.ProxyExport != nil:
			exportFrom, err := l.ResolvePath(stmt.ProxyExport.From, filepath.Dir(file.Path))
			if err != nil {
				errors = append(errors, err)
				continue
			} else if exportFrom == "" {
				continue
			}

			switch {
			case stmt.ProxyExport.ExportAll:
				if stmt.ProxyExport.ExportAllAlias != "" {
					exports = append(exports, language.ExportEntry{
						Names: []language.ExportName{
							{
								Original: stmt.ProxyExport.ExportAllAlias,
							},
						},
						Path: exportFrom,
					})
				} else {
					exports = append(exports, language.ExportEntry{
						All:  true,
						Path: exportFrom,
					})
				}
			case stmt.ProxyExport.ExportDeconstruction != nil:
				names := make([]language.ExportName, 0)
				for _, name := range stmt.ProxyExport.ExportDeconstruction.Names {
					names = append(names, language.ExportName{
						Original: name.Original,
						Alias:    name.Alias,
					})
				}

				exports = append(exports, language.ExportEntry{
					Names: names,
					Path:  exportFrom,
				})
			}
		}
	}
	return &language.ExportsEntries{
		Exports: exports,
		Errors:  errors,
	}, nil
}
