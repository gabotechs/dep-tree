package js

import (
	"fmt"
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/js/js_grammar"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/utils"
)

var Extensions = []string{
	"js", "ts", "tsx", "jsx", "d.ts", "mjs", "cjs",
}

type Language struct {
	Cfg *Config
}

var _ language.Language[js_grammar.File] = &Language{}

// findPackageJson starts from a search path and goes up dir by dir
// until a package.json file is found. If one is found, it returns the
// dir where it was found and a parsed TsConfig object in case that there
// was also a tsconfig.json file.
func _findClosestPackageJsonPath(searchPath string) string {
	packageJsonPath := filepath.Join(searchPath, packageJsonFile)
	if utils.FileExists(packageJsonPath) {
		return packageJsonPath
	}
	nextSearchPath := filepath.Dir(searchPath)
	if nextSearchPath != searchPath {
		return _findClosestPackageJsonPath(nextSearchPath)
	} else {
		return ""
	}
}

var findClosestPackageJsonPath = utils.Cached1In1Out(_findClosestPackageJsonPath)

func MakeJsLanguage(entrypoint string, cfg *Config) (language.Language[js_grammar.File], error) {
	if !utils.FileExists(entrypoint) {
		return nil, fmt.Errorf("file %s does not exist", entrypoint)
	}

	return &Language{Cfg: cfg}, nil
}

func (l *Language) ParseFile(id string) (*js_grammar.File, error) {
	return js_grammar.Parse(id)
}
