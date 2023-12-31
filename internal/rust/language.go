package rust

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/rust/rust_grammar"
	"github.com/gabotechs/dep-tree/internal/utils"
)

var Extensions = []string{
	"rs",
}

type Language struct {
	CargoTomlPath     string
	ProjectEntrypoint string
	ModTree           *ModTree
}

func (l *Language) ParseFile(id string) (*rust_grammar.File, error) {
	return CachedRustFile(id)
}

var _ language.Language[rust_grammar.File] = &Language{}

func findCargoToml(searchPath string) string {
	if len(searchPath) < 2 {
		return ""
	} else if p := filepath.Join(searchPath, "Cargo.toml"); utils.FileExists(p) {
		return p
	}
	return findCargoToml(filepath.Dir(searchPath))
}

func findProjectEntrypoint(rootPath string, searchPaths []string) string {
	for _, searchPath := range searchPaths {
		if p := filepath.Join(rootPath, searchPath); utils.FileExists(p) {
			return p
		}
	}
	return ""
}

var searchPaths = []string{
	filepath.Join("src", "lib.rs"),
	filepath.Join("src", "main.rs"),
	"lib.rs",
	"main.rs",
}

func MakeRustLanguage(entrypoint string, _ *Config) (language.Language[rust_grammar.File], error) {
	entrypointAbsPath, err := filepath.Abs(entrypoint)
	if err != nil {
		return nil, err
	}
	cargoTomlPath := findCargoToml(entrypointAbsPath)
	if cargoTomlPath == "" {
		return nil, fmt.Errorf("could not find Cargo.toml in any parent directory of %s", entrypointAbsPath)
	}

	projectEntrypoint := findProjectEntrypoint(filepath.Dir(cargoTomlPath), searchPaths)
	if projectEntrypoint == "" {
		return nil, fmt.Errorf("could not find any of the possible entrypoint paths %s", strings.Join(searchPaths, ", "))
	}

	modTree, err := MakeModTree(projectEntrypoint, "crate")
	if err != nil {
		return nil, err
	}

	return &Language{
		CargoTomlPath:     cargoTomlPath,
		ProjectEntrypoint: projectEntrypoint,
		ModTree:           modTree,
	}, nil
}
