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

var _ language.Language = &Language{}

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
	if nextSearchPath != searchPath {
		findFirstPackageJsonWithNameCache[searchPath] = findFirstPackageJsonWithName(nextSearchPath)
	}
	return findFirstPackageJsonWithNameCache[searchPath]
}

func MakeJsLanguage(cfg *Config) (language.Language, error) {
	return &Language{Cfg: cfg}, nil
}

func (l *Language) ParseFile(id string) (*language.FileInfo, error) {
	fileInfo, err := js_grammar.Parse(id)
	if err != nil {
		return nil, err
	}
	pkgJson := findFirstPackageJsonWithName(filepath.Dir(id))
	if pkgJson == nil {
		return fileInfo, nil
	}
	fileInfo.Package = pkgJson.Name
	fileInfo.RelPath, _ = filepath.Rel(pkgJson.absPath, id)
	return fileInfo, nil
}
