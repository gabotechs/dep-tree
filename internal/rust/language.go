package rust

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"dep-tree/internal/language"
	"dep-tree/internal/rust/rust_grammar"
	"dep-tree/internal/utils"
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

func (l *Language) lazyLoadModTree(ctx context.Context) (context.Context, error) {
	if l.ModTree == nil {
		ctx, modTree, err := MakeModTree(ctx, l.ProjectEntrypoint, "crate", nil)
		l.ModTree = modTree
		return ctx, err
	}
	return ctx, nil
}

func MakeRustLanguage(entrypoint string) (language.Language[Data, rust_grammar.File], error) {
	entrypointAbsPath, err := filepath.Abs(entrypoint)
	if err != nil {
		return nil, err
	}
	cargoTomlPath := findCargoToml(entrypointAbsPath)
	if cargoTomlPath == "" {
		return nil, fmt.Errorf("could not find Cargo.toml in any parent directory of %s", entrypointAbsPath)
	}

	projectEntrypoint := findProjectEntrypoint(path.Dir(cargoTomlPath), searchPaths)
	if projectEntrypoint == "" {
		return nil, fmt.Errorf("could not find any of the possible entrypoint paths %s", strings.Join(searchPaths, ", "))
	}

	return &Language{
		CargoTomlPath:     cargoTomlPath,
		ProjectEntrypoint: projectEntrypoint,
	}, nil
}
