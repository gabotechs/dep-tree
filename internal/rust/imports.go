package rust

import (
	"fmt"
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/rust/rust_grammar"
	"github.com/gabotechs/dep-tree/internal/utils"
)

func (l *Language) ParseImports(file *language.FileInfo) (*language.ImportsResult, error) {
	imports := make([]language.ImportEntry, 0)
	var errors []error

	content := file.Content.(*rust_grammar.File)
	for _, stmt := range content.Statements {
		if stmt.Use != nil {
			for _, use := range stmt.Use.Flatten() {
				id, err := resolve(use.PathSlices, file.Path)
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
						Names: []string{string(use.Name.Original)},
						Path:  id,
					})
				}
			}
		} else if stmt.Mod != nil && !stmt.Mod.Local {
			names := []string{string(stmt.Mod.Name)}

			thisDir := filepath.Dir(file.Path)

			var modPath string
			if p := filepath.Join(thisDir, string(stmt.Mod.Name)+".rs"); utils.FileExists(p) {
				modPath = p
			} else if p = filepath.Join(thisDir, string(stmt.Mod.Name), "mod.rs"); utils.FileExists(p) {
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
