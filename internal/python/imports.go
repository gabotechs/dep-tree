package python

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"dep-tree/internal/language"
	"dep-tree/internal/python/python_grammar"
)

func (l *Language) ParseImports(ctx context.Context, file *python_grammar.File) (context.Context, *language.ImportsResult, error) {
	imports := make([]language.ImportEntry, 0)
	var errors []error

	for _, stmt := range file.Statements {
		switch {
		case stmt == nil:
			// Is this even possible?
		case stmt.Import != nil && !stmt.Import.Indented:
			resolved, err := l.ResolveAbsolute(stmt.Import.Path[0:])
			if err != nil {
				errors = append(errors, err)
			} else if resolved == nil {
				continue
			} else if resolved.File != "" {
				imports = append(imports, language.ImportEntry{
					All: true,
					Id:  resolved.File,
				})
			} else if resolved.InitModule != "" {
				imports = append(imports, language.ImportEntry{
					All: true,
					Id:  resolved.InitModule,
				})
			} else if resolved.Directory != "" {
				// if import references a dir, it is importing all the files in the dir.
				entries, err := os.ReadDir(resolved.Directory)
				if err != nil {
					errors = append(errors, err)
					continue
				}
				for _, entry := range entries {
					if strings.HasSuffix(entry.Name(), ".py") {
						imports = append(imports, language.ImportEntry{
							All: true,
							Id:  path.Join(resolved.Directory, entry.Name()),
						})
					}
				}
			}
		case stmt.FromImport != nil && !stmt.FromImport.Indented:
			names := make([]string, len(stmt.FromImport.Names))
			for i, name := range stmt.FromImport.Names {
				names[i] = name.Name
			}
			if len(names) == 0 {
				names = nil
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
				imports = append(imports, language.ImportEntry{
					Names: names,
					All:   stmt.FromImport.All,
					Id:    resolved.File,
				})
			} else if resolved.InitModule != "" {
				// If imported from an __init__.py, first look if there are files with that name.
				// if there are, then it is importing that file, not something inside __init__.py.
				directory := path.Dir(resolved.InitModule)
				entries, err := os.ReadDir(directory)
				if err != nil {
					errors = append(errors, err)
					continue
				}

				availableFiles := map[string]bool{}
				for _, entry := range entries {
					if strings.HasSuffix(entry.Name(), ".py") {
						availableFiles[entry.Name()] = true
					}
				}

				var namesFromInit []string
				for _, name := range names {
					namePy := name + ".py"
					if _, ok := availableFiles[namePy]; !ok {
						// No file named that way, it should be imported from __init__.py then.
						namesFromInit = append(namesFromInit, name)
					} else {
						// Imported a specific file.
						imports = append(imports, language.ImportEntry{
							All: true,
							Id:  path.Join(directory, namePy),
						})
					}
				}
				if namesFromInit != nil || stmt.FromImport.All {
					imports = append(imports, language.ImportEntry{
						All:   stmt.FromImport.All,
						Names: namesFromInit,
						Id:    resolved.InitModule,
					})
				}
			} else if resolved.Directory != "" {
				// Imported files in that directory.
				entries, err := os.ReadDir(resolved.Directory)
				if err != nil {
					errors = append(errors, err)
					continue
				}

				availableFiles := map[string]bool{}
				for _, entry := range entries {
					if strings.HasSuffix(entry.Name(), ".py") {
						availableFiles[entry.Name()] = true
					}
				}

				// TODO: what happens if stmt.FromImport.All == true.
				for _, name := range names {
					namePy := name + ".py"
					if _, ok := availableFiles[namePy]; !ok {
						errors = append(errors, fmt.Errorf("cannot import file %s from directory %s", namePy, resolved.Directory))
					} else {
						imports = append(imports, language.ImportEntry{
							All: true,
							Id:  path.Join(resolved.Directory, namePy),
						})
					}
				}
			}
		}
	}
	return ctx, &language.ImportsResult{Imports: imports, Errors: errors}, nil
}
