package rust

import (
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/language"
)

var Extensions = []string{
	"rs",
}

func (l *Language) MakeNode(path string) (*graph.Node[language.CodeFile], error) {
	abs, err := filepath.Abs(path)
	return graph.MakeNode(abs, language.CodeFile{}), err
}
