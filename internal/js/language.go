package js

import (
	"context"
	"fmt"
	"path"

	"github.com/gabotechs/dep-tree/internal/js/js_grammar"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/utils"
)

var Extensions = []string{
	"js", "ts", "tsx", "jsx", "d.ts",
}

type Language struct {
	Workspaces *Workspaces
	Cfg        *Config
}

var _ language.Language[js_grammar.File] = &Language{}

// findPackageJson starts from a search path and goes up dir by dir
// until a package.json file is found. If one is found, it returns the
// dir where it was found and a parsed TsConfig object in case that there
// was also a tsconfig.json file.
func _findPackageJson(searchPath string) (TsConfig, string, error) {
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
		return _findPackageJson(path.Dir(searchPath))
	}
}

var findPackageJson = utils.Cached1In2OutErr(_findPackageJson)

func MakeJsLanguage(ctx context.Context, entrypoint string, cfg *Config) (context.Context, language.Language[js_grammar.File], error) {
	if !utils.FileExists(entrypoint) {
		return ctx, nil, fmt.Errorf("file %s does not exist", entrypoint)
	}
	workspaces, err := NewWorkspaces(entrypoint)
	if err != nil {
		return ctx, nil, err
	}

	return ctx, &Language{
		Cfg:        cfg,
		Workspaces: workspaces,
	}, nil
}

func (l *Language) ParseFile(id string) (*js_grammar.File, error) {
	return js_grammar.Parse(id)
}
