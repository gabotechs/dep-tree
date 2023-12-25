package python

import (
	"path"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python/python_grammar"
)

//nolint:gocyclo
func (l *Language) ParseExports(file *python_grammar.File) (*language.ExportsEntries, error) {
	var exports []language.ExportEntry
	var errors []error
	for _, stmt := range file.Statements {
		switch {
		case stmt == nil:
			continue
		case stmt.Import != nil && !stmt.Import.Indented:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{
					{
						Original: stmt.Import.Path[0],
						Alias:    stmt.Import.Alias,
					},
				},
				Path: file.Path,
			})
		case stmt.FromImport != nil && !stmt.FromImport.Indented:
			entry := language.ExportEntry{
				Names: make([]language.ExportName, len(stmt.FromImport.Names)),
				All:   stmt.FromImport.All,
				Path:  file.Path,
			}
			for i, name := range stmt.FromImport.Names {
				entry.Names[i] = language.ExportName{
					Original: name.Name,
					Alias:    name.Alias,
				}
			}
			var resolved *ResolveResult
			if len(stmt.FromImport.Relative) > 0 {
				var err error
				resolved, err = ResolveRelative(stmt.FromImport.Path, path.Dir(file.Path), len(stmt.FromImport.Relative)-1)
				if err != nil {
					errors = append(errors, err)
					continue
				}
			} else {
				resolved = l.ResolveAbsolute(stmt.FromImport.Path)
			}
			switch {
			case resolved == nil:
			case resolved.InitModule != nil:
			case resolved.Directory != nil:
				// nothing.
			case resolved.File != nil:
				entry.Path = resolved.File.Path
			}
			exports = append(exports, entry)

		case stmt.VariableUnpack != nil:
			entry := language.ExportEntry{
				Names: make([]language.ExportName, len(stmt.VariableUnpack.Names)),
				Path:  file.Path,
			}
			for i, name := range stmt.VariableUnpack.Names {
				entry.Names[i] = language.ExportName{Original: name}
			}
			exports = append(exports, entry)
		case stmt.VariableAssign != nil:
			entry := language.ExportEntry{
				Names: make([]language.ExportName, len(stmt.VariableAssign.Names)),
				Path:  file.Path,
			}
			for i, name := range stmt.VariableAssign.Names {
				entry.Names[i] = language.ExportName{Original: name}
			}
			exports = append(exports, entry)
		case stmt.VariableTyping != nil:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{{Original: stmt.VariableTyping.Name}},
				Path:  file.Path,
			})
		case stmt.Function != nil:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{{Original: stmt.Function.Name}},
				Path:  file.Path,
			})
		case stmt.Class != nil:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{{Original: stmt.Class.Name}},
				Path:  file.Path,
			})
		}
	}
	return &language.ExportsEntries{Exports: exports, Errors: errors}, nil
}
