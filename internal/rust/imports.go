package rust

import (
	"fmt"
	"path"

	"dep-tree/internal/language"
	"dep-tree/internal/rust/rust_grammar"
	"dep-tree/internal/utils"
)

func (l *Language) ParseImports(file *rust_grammar.File) (*language.ImportsResult, error) {
	imports := make([]language.ImportEntry, 0)
	var errors []error

	for _, stmt := range file.Statements {
		if stmt.Use != nil {
			names := make([]string, len(stmt.Use.Names))
			for i, name := range stmt.Use.Names {
				names[i] = name.Original
			}

			id, err := l.resolve(stmt.Use.PathSlices, file.Path)
			if err != nil {
				errors = append(errors, err)
				continue
			} else if id == "" {
				continue
			}

			imports = append(imports, language.ImportEntry{
				All:   stmt.Use.All,
				Names: names,
				Id:    id,
			})
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
				Id:    modPath,
			})
		}
	}

	return &language.ImportsResult{
		Imports: imports,
		Errors:  errors,
	}, nil
}
