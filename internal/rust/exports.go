package rust

import (
	"context"
	"fmt"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/rust/rust_grammar"
)

func (l *Language) ParseExports(ctx context.Context, file *rust_grammar.File) (context.Context, *language.ExportsResult, error) {
	exports := make([]language.ExportEntry, 0)
	var errors []error

	for _, stmt := range file.Statements {
		switch {
		case stmt.Use != nil && stmt.Use.Pub:
			for _, use := range stmt.Use.Flatten() {
				newCtx, id, err := l.resolve(ctx, use.PathSlices, file.Path)
				ctx = newCtx
				if err != nil {
					errors = append(errors, fmt.Errorf("error resolving use statement for name %s: %w", use.Name.Original, err))
					continue
				} else if id == "" {
					continue
				}

				if use.All {
					exports = append(exports, language.ExportEntry{
						All: use.All,
						Id:  id,
					})
				} else {
					exports = append(exports, language.ExportEntry{
						Names: []language.ExportName{{Original: use.Name.Original, Alias: use.Name.Alias}},
						Id:    id,
					})
				}
			}
		case stmt.Pub != nil:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{{Original: stmt.Pub.Name}},
				Id:    file.Path,
			})
		case stmt.Mod != nil && stmt.Mod.Pub:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{{Original: stmt.Mod.Name}},
				Id:    file.Path,
			})
		}
	}

	return ctx, &language.ExportsResult{
		Exports: exports,
		Errors:  errors,
	}, nil
}
