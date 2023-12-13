package rust

import (
	"fmt"
	"path"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/rust/rust_grammar"
	"github.com/gabotechs/dep-tree/internal/utils"
)

func (l *Language) ParseImports(file *rust_grammar.File) (*language.ImportsResult, error) {
	imports := make([]language.ImportEntry, 0)
	var errors []error

	for _, stmt := range file.Statements {
		if stmt.Use != nil {
			for _, use := range stmt.Use.Flatten() {
				id, err := l.resolve(use.PathSlices, file.Path)
				if err != nil {
					errors = append(errors, fmt.Errorf("error resolving use statement for name %s: %w", use.Name.Original, err))
					continue
				} else if id == "" {
					continue
				}

				if use.All {
					imports = append(imports, language.ImportEntry{
						All:  use.All,
						Path: id,
					})
				} else {
					imports = append(imports, language.ImportEntry{
						Names: []string{use.Name.Original},
						Path:  id,
					})
				}
			}
		} else if stmt.Mod != nil && !stmt.Mod.Local {
			names := []string{stmt.Mod.Name}

			thisDir := path.Dir(file.Path)

			var modPath string
			if p := path.Join(thisDir, stmt.Mod.Name+".rs"); utils.FileExists(p) {
				modPath = p
			} else if p = path.Join(thisDir, stmt.Mod.Name, "mod.rs"); utils.FileExists(p) {
				modPath = p
			} else {
				errors = append(errors, fmt.Errorf("could not find mod %s while looking in dir %s", stmt.Mod.Name, thisDir))
				continue
			}

			imports = append(imports, language.ImportEntry{
				All:   true,
				Names: names,
				Path:  modPath,
			})
		}
	}

	return &language.ImportsResult{
		Imports: imports,
		Errors:  errors,
	}, nil
}
