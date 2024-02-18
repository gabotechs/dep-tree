package js

import (
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/js/js_grammar"
	"github.com/gabotechs/dep-tree/internal/language"
)

type ExportsCacheKey string

func (l *Language) ParseExports(file *language.FileInfo) (*language.ExportsResult, error) {
	exports := make([]language.ExportEntry, 0)
	var errors []error

	content := file.Content.(*js_grammar.File)
	for _, stmt := range content.Statements {
		switch {
		case stmt == nil:
			// Is this even possible?
		case stmt.DeclarationExport != nil:
			exports = append(exports, language.ExportEntry{
				Symbols: []language.ExportSymbol{
					{
						Original: stmt.DeclarationExport.Name,
					},
				},
				AbsPath: file.AbsPath,
			})
		case stmt.ListExport != nil:
			if stmt.ListExport.ExportDeconstruction != nil {
				for _, name := range stmt.ListExport.ExportDeconstruction.Names {
					exports = append(exports, language.ExportEntry{
						Symbols: []language.ExportSymbol{
							{
								Original: name.Original,
								Alias:    name.Alias,
							},
						},
						AbsPath: file.AbsPath,
					})
				}
			}
		case stmt.DefaultExport != nil:
			if stmt.DefaultExport.Default {
				exports = append(exports, language.ExportEntry{
					Symbols: []language.ExportSymbol{
						{
							Original: "default",
						},
					},
					AbsPath: file.AbsPath,
				})
			}
		case stmt.ProxyExport != nil:
			exportFrom, err := l.ResolvePath(stmt.ProxyExport.From, filepath.Dir(file.AbsPath))
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
						Symbols: []language.ExportSymbol{
							{
								Original: stmt.ProxyExport.ExportAllAlias,
							},
						},
						AbsPath: exportFrom,
					})
				} else {
					exports = append(exports, language.ExportEntry{
						All:     true,
						AbsPath: exportFrom,
					})
				}
			case stmt.ProxyExport.ExportDeconstruction != nil:
				names := make([]language.ExportSymbol, 0)
				for _, name := range stmt.ProxyExport.ExportDeconstruction.Names {
					names = append(names, language.ExportSymbol{
						Original: name.Original,
						Alias:    name.Alias,
					})
				}

				exports = append(exports, language.ExportEntry{
					Symbols: names,
					AbsPath: exportFrom,
				})
			}
		}
	}
	return &language.ExportsResult{
		Exports: exports,
		Errors:  errors,
	}, nil
}
