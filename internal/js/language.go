package js

import (
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

var findFirstPackageJsonWithNameCache = map[string]*packageJson{}

func findFirstPackageJsonWithName(searchPath string) *packageJson {
	if result, ok := findFirstPackageJsonWithNameCache[searchPath]; ok {
		return result
	}
	packageJsonPath := filepath.Join(searchPath, packageJsonFile)
	if utils.FileExists(packageJsonPath) {
		pckJson, _ := readPackageJson(packageJsonPath)
		if pckJson != nil && pckJson.Name != "" {
			return pckJson
		}
	}
	nextSearchPath := filepath.Dir(searchPath)
	findFirstPackageJsonWithNameCache[searchPath] = findFirstPackageJsonWithName(nextSearchPath)
	return findFirstPackageJsonWithNameCache[searchPath]
}

func (l *Language) Display(id string) language.DisplayResult {
	pkgJson := findFirstPackageJsonWithName(filepath.Dir(id))
	if pkgJson == nil {
		return language.DisplayResult{Name: id}
	}

	result, err := filepath.Rel(pkgJson.absPath, id)
	if err != nil {
		return language.DisplayResult{Name: id, Group: pkgJson.Name}
	}
	return language.DisplayResult{Name: result, Group: pkgJson.Name}
}

func MakeJsLanguage(cfg *Config) (language.Language[js_grammar.File], error) {
	return &Language{Cfg: cfg}, nil
}

func (l *Language) ParseFile(id string) (*js_grammar.File, error) {
	return js_grammar.Parse(id)
}
