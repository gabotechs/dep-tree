package js

import (
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

func MakeJsLanguage(cfg *Config) (language.Language[js_grammar.File], error) {
	return &Language{Cfg: cfg}, nil
}

func (l *Language) ParseFile(id string) (*js_grammar.File, error) {
	return js_grammar.Parse(id)
}
