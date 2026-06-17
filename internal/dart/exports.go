package dart

import (
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/dart/dart_grammar"
	"github.com/gabotechs/dep-tree/internal/language"
)

func (l *Language) ParseExports(file *language.FileInfo) (*language.ExportsResult, error) {
	exports := make([]language.ExportEntry, 0)
	var errors []error

	dir := filepath.Dir(file.AbsPath)
	content := file.Content.(*dart_grammar.File)
	for _, stmt := range content.Statements {
		switch {
		case stmt == nil:
			continue

		// An `export` directive re-exports symbols declared in another file.
		case stmt.Export != nil:
			exportFrom, err := l.ResolvePath(stmt.Export.Path, dir)
			if err != nil {
				errors = append(errors, err)
				continue
			} else if exportFrom == "" {
				continue
			}

			var shown []string
			for _, combinator := range stmt.Export.Combinators {
				shown = append(shown, combinator.Show...)
			}
			if len(shown) > 0 {
				symbols := make([]language.ExportSymbol, len(shown))
				for i, name := range shown {
					symbols[i] = language.ExportSymbol{Original: name}
				}
				exports = append(exports, language.ExportEntry{
					Symbols: symbols,
					AbsPath: exportFrom,
				})
			} else {
				exports = append(exports, language.ExportEntry{
					All:     true,
					AbsPath: exportFrom,
				})
			}

		// A top-level declaration exposes a public symbol declared in this same file.
		case stmt.Declaration != nil:
			exports = append(exports, language.ExportEntry{
				Symbols: []language.ExportSymbol{{Original: stmt.Declaration.Name}},
				AbsPath: file.AbsPath,
			})
		}
	}

	return &language.ExportsResult{
		Exports: exports,
		Errors:  errors,
	}, nil
}
