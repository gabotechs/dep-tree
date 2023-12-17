package js

import (
	"context"
	"fmt"
	"path"
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/js/js_grammar"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/utils"
)

type Language struct {
	PackageJsonPath string
	ProjectRoot     string
	TsConfig        TsConfig
	Cfg             *Config
}

var _ language.Language[js_grammar.File] = &Language{}

func findPackageJson(searchPath string) (TsConfig, string, error) {
	if len(searchPath) < 2 {
		return TsConfig{}, "", nil
	}
	packageJsonPath := path.Join(searchPath, "package.json")
	if utils.FileExists(packageJsonPath) {
		tsConfigPath := path.Join(searchPath, "tsconfig.json")
		var tsConfig TsConfig
		var err error
		if utils.FileExists(tsConfigPath) {
			tsConfig, err = ParseTsConfig(tsConfigPath)
			if err != nil {
				err = fmt.Errorf("found TypeScript config file in %s but there was an error reading it: %w", tsConfigPath, err)
			}
		}
		return tsConfig, searchPath, err
	} else {
		return findPackageJson(path.Dir(searchPath))
	}
}

func MakeJsLanguage(ctx context.Context, entrypoint string, cfg *Config) (context.Context, language.Language[js_grammar.File], error) {
	entrypointAbsPath, err := filepath.Abs(entrypoint)
	if err != nil {
		return ctx, nil, err
	}
	if !utils.FileExists(entrypoint) {
		return ctx, nil, fmt.Errorf("file %s does not exist", entrypoint)
	}

	tsConfig, packageJsonPath, err := findPackageJson(entrypointAbsPath)
	if err != nil {
		return ctx, nil, err
	}
	projectRoot := path.Dir(entrypointAbsPath)
	if packageJsonPath != "" {
		projectRoot = packageJsonPath
	}
	return ctx, &Language{
		PackageJsonPath: packageJsonPath,
		ProjectRoot:     projectRoot,
		TsConfig:        tsConfig,
		Cfg:             cfg,
	}, nil
}

func (l *Language) ParseFile(id string) (*js_grammar.File, error) {
	return js_grammar.Parse(id)
}
