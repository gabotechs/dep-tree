package rust

import (
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/graph"
)

var Extensions = []string{
	"rs",
}

type Data struct{}

func (l *Language) MakeNode(path string) (*graph.Node[Data], error) {
	abs, err := filepath.Abs(path)
	return graph.MakeNode(abs, Data{}), err
}
