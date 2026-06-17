package dart

import (
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/dart/dart_grammar"
	"github.com/gabotechs/dep-tree/internal/language"
)

var Extensions = []string{
	"dart",
}

type Language struct {
	Cfg *Config
}

var _ language.Language = &Language{}

func MakeDartLanguage(cfg *Config) (language.Language, error) {
	if cfg == nil {
		cfg = &Config{}
	}
	return &Language{Cfg: cfg}, nil
}

func (l *Language) ParseFile(id string) (*language.FileInfo, error) {
	fileInfo, err := dart_grammar.Parse(id)
	if err != nil {
		return nil, err
	}
	// The "root of the project" for a Dart file is the directory that holds the closest
	// pubspec.yaml file. The package name is taken from the `name` attribute of that file.
	pubspec := findClosestPubspec(filepath.Dir(id))
	if pubspec == nil {
		return fileInfo, nil
	}
	fileInfo.Package = pubspec.Name
	fileInfo.RelPath, _ = filepath.Rel(pubspec.absDir, id)
	return fileInfo, nil
}
