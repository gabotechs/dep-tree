package python

import (
	"fmt"
	"path"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python/python_grammar"
)

// handleImport Handles an `import` statement, e.g. `import foo`
func (l *Language) handleImport(imp *python_grammar.Import) []language.ImportEntry {
	if l.cfg.ExcludeConditionalImports && imp.Indented {
		return nil
	}
	resolved := l.ResolveAbsolute(imp.Path[0:])
	if resolved == nil {
		return nil
	}
	switch {
	// `import my_file` -> all names from my_file.py are imported.
	case resolved.File != nil:
		return []language.ImportEntry{language.EmptyImport(resolved.File.Path)}

	// `import my_module` -> all names from my_module/__init__.py are imported.
	case resolved.InitModule != nil:
		return []language.ImportEntry{language.EmptyImport(resolved.InitModule.Path)}

	// `import my_dir` -> all python files inside dir my_dir/*.py are imported, but no specific name.
	case resolved.Directory != nil && !l.IgnoreDirectoryImports:
		imports := make([]language.ImportEntry, len(resolved.Directory.PythonFiles))
		for i, pythonFile := range resolved.Directory.PythonFiles {
			imports[i] = language.EmptyImport(pythonFile)
		}
		return imports
	}

	return nil
}

// handleFromImportAll handles a `from X import *` statement, e.g. `from foo import *`
func handleFromImportAll(resolved *ResolveResult) []language.ImportEntry {
	switch {
	// `from my_file import *` -> all names from my_file.py are imported.
	case resolved.File != nil:
		return []language.ImportEntry{language.AllImport(resolved.File.Path)}

	// `from my_module import *` -> all names from the my_module/__init__.py file are imported.
	case resolved.InitModule != nil:
		return []language.ImportEntry{language.AllImport(resolved.InitModule.Path)}

	// `from my_folder import *` -> nothing is imported.
	case resolved.Directory != nil:
		return nil
	}
	return nil
}

// handleFromImportNames handles a `from X import a,b,c` statement, e.g. `from foo import bar, baz`
func handleFromImportNames(resolved *ResolveResult, names []string) ([]language.ImportEntry, error) {
	switch {
	// `from my_file import foo, bar` -> names foo and bar from my_file.py are imported.
	case resolved.File != nil:
		return []language.ImportEntry{language.NamesImport(names, resolved.File.Path)}, nil

	// `from my_module import foo, bar` -> names foo and bar from the my_module/__init__.py file are imported.
	// It might happen that some of those names are actually Python files (e.g. my_module/foo.py).
	// TODO: If a name first matches a Python file then assume that the file is being imported, otherwise,
	//  it must be imported from my_module/__init__.py. This is wrong, but here we do not have any info about
	//  what names are exported from my_module/__init__.py, so the best we can do is just first try to check
	//  for existing Python files.
	case resolved.InitModule != nil:
		availableFiles := resolved.InitModule.fileMap()

		var imports []language.ImportEntry
		var namesFromInit []string
		for _, name := range names {
			if pythonFile, ok := availableFiles[name]; ok {
				// Imported a specific file.
				imports = append(imports, language.EmptyImport(pythonFile))
			} else {
				// No file named that way, it should be imported from __init__.py then.
				namesFromInit = append(namesFromInit, name)
			}
		}
		if namesFromInit != nil {
			imports = append(imports, language.NamesImport(namesFromInit, resolved.InitModule.Path))
		}
		return imports, nil

	// `from my_folder import *` -> all the files are imported.
	case resolved.Directory != nil:
		availableFiles := resolved.Directory.fileMap()
		for _, name := range names {
			if pythonFile, ok := availableFiles[name]; ok {
				return []language.ImportEntry{language.EmptyImport(pythonFile)}, nil
			} else {
				return nil, fmt.Errorf(
					"cannot import file %s.py from directory %s",
					name,
					resolved.Directory.Path,
				)
			}
		}
	}
	return nil, nil
}

func (l *Language) handleFromImport(imp *python_grammar.FromImport, currDir string) ([]language.ImportEntry, error) {
	if l.cfg.ExcludeConditionalImports && imp.Indented {
		return nil, nil
	}

	resolved, err := l.resolveFromImportPath(imp, currDir)
	if err != nil || resolved == nil {
		return nil, err
	}

	if imp.All {
		return handleFromImportAll(resolved), nil
	} else {
		names := make([]string, len(imp.Names))
		for i, name := range imp.Names {
			names[i] = name.Name
		}

		return handleFromImportNames(resolved, names)
	}
}

func (l *Language) ParseImports(file *python_grammar.File) (*language.ImportsResult, error) {
	imports := make([]language.ImportEntry, 0)
	var errors []error

	for _, stmt := range file.Statements {
		switch {
		case stmt == nil:
			// Is this even possible?
		case stmt.Import != nil:
			imports = append(imports, l.handleImport(stmt.Import)...)
		case stmt.FromImport != nil:
			newImports, err := l.handleFromImport(stmt.FromImport, path.Dir(file.Path))
			imports = append(imports, newImports...)
			if err != nil {
				errors = append(errors, err)
			}
		}
	}
	return &language.ImportsResult{Imports: imports, Errors: errors}, nil
}

func (l *Language) resolveFromImportPath(imp *python_grammar.FromImport, currDir string) (*ResolveResult, error) {
	if len(imp.Relative) > 0 {
		return ResolveRelative(imp.Path, currDir, len(imp.Relative)-1)
	} else {
		return l.ResolveAbsolute(imp.Path), nil
	}
}
