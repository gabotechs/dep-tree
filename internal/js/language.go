package js

import (
	"fmt"
	"path"
	"path/filepath"

	"dep-tree/internal/js/js_grammar"
	"dep-tree/internal/language"
	"dep-tree/internal/utils"
)

type Language struct {
	PackageJsonPath string
	ProjectRoot     string
	TsConfig        TsConfig
}

var _ language.Language[Data, js_grammar.File] = &Language{}

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

func MakeJsLanguage(entrypoint string) (language.Language[Data, js_grammar.File], error) {
	entrypointAbsPath, err := filepath.Abs(entrypoint)
	if err != nil {
		return nil, err
	}
	tsConfig, packageJsonPath, err := findPackageJson(entrypointAbsPath)
	if err != nil {
		return nil, err
	}
	projectRoot := path.Dir(entrypointAbsPath)
	if packageJsonPath != "" {
		projectRoot = packageJsonPath
	}
	return &Language{
		PackageJsonPath: packageJsonPath,
		ProjectRoot:     projectRoot,
		TsConfig:        tsConfig,
	}, nil
}

func (l *Language) ParseFile(id string) (*js_grammar.File, error) {
	return js_grammar.Parse(id)
}
