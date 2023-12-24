package python

import (
	"fmt"
	"path"
	"strings"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python/python_grammar"
)

//nolint:gocyclo
func (l *Language) ParseImports(file *python_grammar.File) (*language.ImportsResult, error) {
	imports := make([]language.ImportEntry, 0)
	var errors []error

	for _, stmt := range file.Statements {
		switch {
		case stmt == nil:
			// Is this even possible?
		case stmt.Import != nil:
			if l.cfg.ExcludeConditionalImports && stmt.Import.Indented {
				continue
			}
			resolved := l.ResolveAbsolute(stmt.Import.Path[0:])
			switch {
			case resolved == nil:
				continue
			case resolved.File != nil:
				imports = append(imports, language.ImportEntry{
					All:  true,
					Path: resolved.File.Path,
				})
			case resolved.InitModule != nil && !l.IgnoreModuleImports:
				imports = append(imports, language.ImportEntry{
					All:  true,
					Path: resolved.InitModule.Path,
				})
			case resolved.Directory != nil && !l.IgnoreModuleImports:
				for _, pythonFile := range resolved.Directory.PythonFiles {
					imports = append(imports, language.ImportEntry{
						All:  true,
						Path: pythonFile,
					})
				}
			}
		case stmt.FromImport != nil:
			if l.cfg.ExcludeConditionalImports && stmt.FromImport.Indented {
				continue
			}
			importedNames := make([]string, len(stmt.FromImport.Names))
			for i, name := range stmt.FromImport.Names {
				importedNames[i] = name.Name
			}
			if len(importedNames) == 0 {
				importedNames = nil
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
				continue
			case resolved.File != nil:
				imports = append(imports, language.ImportEntry{
					Names: importedNames,
					All:   stmt.FromImport.All,
					Path:  resolved.File.Path,
				})
			case resolved.InitModule != nil:
				// If importing from an __init__.py, there is a chance that we are actually
				// importing a file living in the same folder, instead of a variable that lives
				// inside __init__.py.
				availableFiles := map[string]string{}
				for _, pythonFile := range resolved.InitModule.PythonFiles {
					availableFiles[strings.TrimSuffix(path.Base(pythonFile), ".py")] = pythonFile
				}
				var namesFromInit []string
				for _, name := range importedNames {
					if pythonFile, ok := availableFiles[name]; ok {
						// Imported a specific file.
						imports = append(imports, language.ImportEntry{
							All:  true,
							Path: pythonFile,
						})
					} else {
						// No file named that way, it should be imported from __init__.py then.
						namesFromInit = append(namesFromInit, name)
					}
				}
				if namesFromInit != nil || stmt.FromImport.All {
					imports = append(imports, language.ImportEntry{
						All:   stmt.FromImport.All,
						Names: namesFromInit,
						Path:  resolved.InitModule.Path,
					})
				}
			case resolved.Directory != nil:
				availableFiles := map[string]string{}
				for _, pythonFile := range resolved.Directory.PythonFiles {
					availableFiles[strings.TrimSuffix(path.Base(pythonFile), ".py")] = pythonFile
				}
				for _, name := range importedNames {
					if pythonFile, ok := availableFiles[name]; ok {
						imports = append(imports, language.ImportEntry{
							All:  true,
							Path: pythonFile,
						})
					} else {
						errors = append(
							errors,
							fmt.Errorf(
								"cannot import file %s.py from directory %s",
								name,
								resolved.Directory.Path,
							),
						)
					}
				}
			}
		}
	}
	return &language.ImportsResult{Imports: imports, Errors: errors}, nil
}
