package rust

import (
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/rust/rust_grammar"
)

var Extensions = []string{
	"rs",
}

type Language struct{}

func (l *Language) ParseFile(id string) (*rust_grammar.File, error) {
	return CachedRustFile(id)
}

func (l *Language) Display(id string) language.DisplayResult {
	cargoToml, err := findClosestCargoToml(filepath.Dir(id))
	if err != nil {
		return language.DisplayResult{
			Name: id,
		}
	}
	result, err := filepath.Rel(cargoToml.path, id)
	if err != nil {
		return language.DisplayResult{Name: id, Group: cargoToml.PackageDefinition.Name}
	}
	return language.DisplayResult{Name: result, Group: cargoToml.PackageDefinition.Name}
}

var _ language.Language[rust_grammar.File] = &Language{}

func MakeRustLanguage(_ *Config) (language.Language[rust_grammar.File], error) {
	return &Language{}, nil
}
