package rust

import (
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/language"
)

var Extensions = []string{
	"rs",
}

type Language struct{}

func (l *Language) ParseFile(id string) (*language.FileInfo, error) {
	return CachedRustFile(id)
}

func (l *Language) Display(id string) graph.DisplayResult {
	cargoToml, err := findClosestCargoToml(filepath.Dir(id))
	if err != nil {
		return graph.DisplayResult{
			Name: id,
		}
	}
	result, err := filepath.Rel(cargoToml.path, id)
	if err != nil {
		return graph.DisplayResult{Name: id, Group: cargoToml.PackageDefinition.Name}
	}
	return graph.DisplayResult{Name: result, Group: cargoToml.PackageDefinition.Name}
}

var _ language.Language = &Language{}

func MakeRustLanguage(_ *Config) (language.Language, error) {
	return &Language{}, nil
}
