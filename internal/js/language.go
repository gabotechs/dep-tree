package js

import (
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/js/js_grammar"
	"github.com/gabotechs/dep-tree/internal/language"
)

var Extensions = []string{
	"js", "ts", "tsx", "jsx", "d.ts", "mjs", "cjs",
}

type Language struct {
	Cfg *Config
}

var _ language.Language[js_grammar.File] = &Language{}

func (l *Language) Display(id string) string {
	basePath := filepath.Dir(findClosestPackageJsonPath(filepath.Dir(id)))
	result, err := filepath.Rel(basePath, id)
	if err != nil {
		return id
	}
	return result
}

func MakeJsLanguage(cfg *Config) (language.Language[js_grammar.File], error) {
	return &Language{Cfg: cfg}, nil
}

func (l *Language) ParseFile(id string) (*js_grammar.File, error) {
	return js_grammar.Parse(id)
}
