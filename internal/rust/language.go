package rust

import (
	"fmt"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/rust/rust_grammar"
	"path/filepath"
)

var Extensions = []string{
	"rs",
}

type Language struct{}

func (l *Language) ParseFile(id string) (*rust_grammar.File, error) {
	return CachedRustFile(id)
}

var _ language.Language[rust_grammar.File] = &Language{}

func MakeRustLanguage(entrypoint string, _ *Config) (language.Language[rust_grammar.File], error) {
	entrypointAbsPath, err := filepath.Abs(entrypoint)
	if err != nil {
		return nil, err
	}
	cargoToml, err := findClosestCargoToml(entrypointAbsPath)
	if err != nil {
		return nil, err
	}
	if cargoToml == nil {
		return nil, fmt.Errorf("could not find Cargo.toml in any parent directory of %s", entrypointAbsPath)
	}
	if _, err = cargoToml.MainFile(); err != nil {
		return nil, err
	}

	return &Language{}, nil
}
