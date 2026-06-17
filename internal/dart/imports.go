package dart

import (
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/dart/dart_grammar"
	"github.com/gabotechs/dep-tree/internal/language"
)

func (l *Language) ParseImports(file *language.FileInfo) (*language.ImportsResult, error) {
	imports := make([]language.ImportEntry, 0)
	var errors []error

	dir := filepath.Dir(file.AbsPath)
	content := file.Content.(*dart_grammar.File)
	for _, stmt := range content.Statements {
		var importPath string
		entry := language.ImportEntry{}

		switch {
		case stmt == nil:
			continue
		case stmt.Import != nil:
			importPath = stmt.Import.Path
			fillCombinators(&entry, stmt.Import.Combinators)
		case stmt.Export != nil:
			// An `export` directive also pulls in the exported file as a dependency.
			importPath = stmt.Export.Path
			fillCombinators(&entry, stmt.Export.Combinators)
		case stmt.Part != nil:
			// A `part` directive includes another file into the current library.
			importPath = stmt.Part.Path
			entry.All = true
		default:
			continue
		}

		absPath, err := l.ResolvePath(importPath, dir)
		if err != nil {
			errors = append(errors, err)
		} else if absPath != "" {
			entry.AbsPath = absPath
			imports = append(imports, entry)
		}
	}

	return &language.ImportsResult{
		Imports: imports,
		Errors:  errors,
	}, nil
}

// fillCombinators translates the `show`/`hide` clauses of a directive into the symbols of an
// import entry. A `show` clause imports only the listed symbols, while any other case
// (a `hide` clause or no clause at all) imports all the symbols.
func fillCombinators(entry *language.ImportEntry, combinators []*dart_grammar.Combinator) {
	var shown []string
	for _, combinator := range combinators {
		shown = append(shown, combinator.Show...)
	}
	if len(shown) > 0 {
		entry.Symbols = shown
	} else {
		entry.All = true
	}
}
