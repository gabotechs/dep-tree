package js

import (
	"dep-tree/internal/graph"
)

var Extensions = []string{
	"js", "ts", "tsx", "jsx", "d.ts",
}

type Data struct{}

func MakeJsNode(absFilePath string) (*graph.Node[Data], error) {
	return graph.MakeNode(absFilePath, Data{}), nil
}
