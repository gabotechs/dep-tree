package python

import (
	"path/filepath"

	"dep-tree/internal/graph"
)

var Extensions = []string{
	"py",
}

type Data struct{}

func (l *Language) MakeNode(path string) (*graph.Node[Data], error) {
	abs, err := filepath.Abs(path)
	return graph.MakeNode(abs, Data{}), err
}
