package python

import (
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python/python_grammar"
)

func (l *Language) handleFromImportForExport(imp *python_grammar.FromImport, filePath string) ([]language.ExportEntry, error) {
	resolved, err := l.resolveFromImportPath(imp, filepath.Dir(filePath))
	if err != nil {
		return nil, err
	}

	entry := language.ExportEntry{
		All:  imp.All,
		Path: filePath,
	}
	for _, name := range imp.Names {
		entry.Names = append(entry.Names, language.ExportName{
			Original: name.Name,
			Alias:    name.Alias,
		})
	}
	switch {
	case resolved == nil:
	case resolved.Directory != nil:
	case resolved.InitModule != nil:
		// nothing.
	case resolved.File != nil:
		entry.Path = resolved.File.Path
	}

	return []language.ExportEntry{entry}, nil
}

//nolint:gocyclo
func (l *Language) ParseExports(file *language.FileInfo) (*language.ExportsEntries, error) {
	var exports []language.ExportEntry
	var errors []error

	content := file.Content.(*python_grammar.File)
	for _, stmt := range content.Statements {
		switch {
		case stmt == nil:
			continue
		case stmt.Import != nil && !stmt.Import.Indented && !l.cfg.IgnoreFromImportsAsExports:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{
					{
						Original: stmt.Import.Path[len(stmt.Import.Path)-1],
						Alias:    stmt.Import.Alias,
					},
				},
				Path: file.AbsPath,
			})
		case stmt.FromImport != nil && !stmt.FromImport.Indented && !l.cfg.IgnoreFromImportsAsExports:
			newExports, err := l.handleFromImportForExport(stmt.FromImport, file.AbsPath)
			if err != nil {
				errors = append(errors, err)
			} else {
				exports = append(exports, newExports...)
			}

		case stmt.VariableUnpack != nil:
			entry := language.ExportEntry{
				Names: make([]language.ExportName, len(stmt.VariableUnpack.Names)),
				Path:  file.AbsPath,
			}
			for i, name := range stmt.VariableUnpack.Names {
				entry.Names[i] = language.ExportName{Original: name}
			}
			exports = append(exports, entry)
		case stmt.VariableAssign != nil:
			entry := language.ExportEntry{
				Names: make([]language.ExportName, len(stmt.VariableAssign.Names)),
				Path:  file.AbsPath,
			}
			for i, name := range stmt.VariableAssign.Names {
				entry.Names[i] = language.ExportName{Original: name}
			}
			exports = append(exports, entry)
		case stmt.VariableTyping != nil:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{{Original: stmt.VariableTyping.Name}},
				Path:  file.AbsPath,
			})
		case stmt.Function != nil:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{{Original: stmt.Function.Name}},
				Path:  file.AbsPath,
			})
		case stmt.Class != nil:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{{Original: stmt.Class.Name}},
				Path:  file.AbsPath,
			})
		}
	}
	return &language.ExportsEntries{Exports: exports, Errors: errors}, nil
}
