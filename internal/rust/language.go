package rust

import (
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/language"
)

var Extensions = []string{
	"rs",
}

type Language struct{}

var _ language.Language = &Language{}

func MakeRustLanguage(_ *Config) (language.Language, error) {
	return &Language{}, nil
}

func (l *Language) ParseFile(id string) (*language.FileInfo, error) {
	file, err := CachedRustFile(id)
	if err != nil {
		return nil, err
	}
	cargoToml, err := findClosestCargoToml(filepath.Dir(id))
	if err != nil {
		return file, nil
	}
	file.Package = cargoToml.PackageDefinition.Name
	file.RelPath, _ = filepath.Rel(cargoToml.path, id)
	return file, nil
}
