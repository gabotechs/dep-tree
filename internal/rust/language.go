package rust

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/rust/rust_grammar"
	"github.com/gabotechs/dep-tree/internal/utils"
)

type Language struct {
	CargoTomlPath     string
	ProjectEntrypoint string
	ModTree           *ModTree
}

func (l *Language) ParseFile(id string) (*rust_grammar.File, error) {
	return rust_grammar.Parse(id)
}

var _ language.Language[Data, rust_grammar.File] = &Language{}

func findCargoToml(searchPath string) string {
	if len(searchPath) < 2 {
		return ""
	} else if p := path.Join(searchPath, "Cargo.toml"); utils.FileExists(p) {
		return p
	}
	return findCargoToml(path.Dir(searchPath))
}

func findProjectEntrypoint(rootPath string, searchPaths []string) string {
	for _, searchPath := range searchPaths {
		if p := path.Join(rootPath, searchPath); utils.FileExists(p) {
			return p
		}
	}
	return ""
}

var searchPaths = []string{
	path.Join("src", "lib.rs"),
	path.Join("src", "main.rs"),
	"lib.rs",
	"main.rs",
}

func MakeRustLanguage(ctx context.Context, entrypoint string, _ *Config) (context.Context, language.Language[Data, rust_grammar.File], error) {
	entrypointAbsPath, err := filepath.Abs(entrypoint)
	if err != nil {
		return ctx, nil, err
	}
	cargoTomlPath := findCargoToml(entrypointAbsPath)
	if cargoTomlPath == "" {
		return ctx, nil, fmt.Errorf("could not find Cargo.toml in any parent directory of %s", entrypointAbsPath)
	}

	projectEntrypoint := findProjectEntrypoint(path.Dir(cargoTomlPath), searchPaths)
	if projectEntrypoint == "" {
		return ctx, nil, fmt.Errorf("could not find any of the possible entrypoint paths %s", strings.Join(searchPaths, ", "))
	}

	ctx, modTree, err := MakeModTree(ctx, projectEntrypoint, "crate", nil)
	if err != nil {
		return ctx, nil, err
	}

	return ctx, &Language{
		CargoTomlPath:     cargoTomlPath,
		ProjectEntrypoint: projectEntrypoint,
		ModTree:           modTree,
	}, nil
}
