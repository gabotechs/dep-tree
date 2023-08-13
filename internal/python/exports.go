package python

import (
	"context"
	"path"

	"dep-tree/internal/language"
	"dep-tree/internal/python/python_grammar"
)

func (l *Language) ParseExports(ctx context.Context, file *python_grammar.File) (context.Context, *language.ExportsResult, error) {
	var exports []language.ExportEntry
	var errors []error
	for _, stmt := range file.Statements {
		switch {
		case stmt == nil:
			// Is this even possible?
		case stmt.Import != nil && !stmt.Import.Indented:
			resolved, err := l.ResolveAbsolute(stmt.Import.Path)
			switch {
			case err != nil:
				errors = append(errors, err)
			case resolved == nil:
				// nothing here.
			default:
				exports = append(exports, language.ExportEntry{
					Names: []language.ExportName{
						{
							Original: stmt.Import.Path[0],
							Alias:    stmt.Import.Alias,
						},
					},
					Id: file.Path,
				})
			}
		case stmt.FromImport != nil && !stmt.FromImport.Indented:
			entry := language.ExportEntry{
				Names: make([]language.ExportName, len(stmt.FromImport.Names)),
				All:   stmt.FromImport.All,
			}
			for i, name := range stmt.FromImport.Names {
				entry.Names[i] = language.ExportName{
					Original: name.Name,
					Alias:    name.Alias,
				}
			}
			var resolved *ResolveResult
			var err error
			if len(stmt.FromImport.Relative) > 0 {
				resolved, err = ResolveRelative(stmt.FromImport.Path, path.Dir(file.Path), len(stmt.FromImport.Relative)-1)
			} else {
				resolved, err = l.ResolveAbsolute(stmt.FromImport.Path)
			}

			if err != nil {
				errors = append(errors, err)
				continue
			} else if resolved == nil {
				continue
			} else if resolved.File != "" {
				entry.Id = resolved.File
			} else if resolved.InitModule != "" {
				// If set to `resolved.InitModule`, it will lead to potential circular exports.
				entry.Id = file.Path
			} else if resolved.Directory != "" {
				// we are exporting the files themselves.
				entry.Id = file.Path
			}
			exports = append(exports, entry)

		case stmt.Variable != nil && !stmt.Variable.Indented:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{
					{
						Original: stmt.Variable.Name,
					},
				},
				Id: file.Path,
			})
		case stmt.Function != nil && !stmt.Function.Indented:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{
					{
						Original: stmt.Function.Name,
					},
				},
				Id: file.Path,
			})
		case stmt.Class != nil && !stmt.Class.Indented:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{
					{
						Original: stmt.Class.Name,
					},
				},
				Id: file.Path,
			})
		}
	}
	return ctx, &language.ExportsResult{Exports: exports, Errors: errors}, nil
}
