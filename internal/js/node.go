package js

import (
	"dep-tree/internal/graph/node"
)

var Extensions = []string{
	"js", "ts", "tsx", "jsx", "d.ts",
}

type Data struct {
	filePath string
}

func MakeJsNode(absFilePath string) (*node.Node[Data], error) {
	return node.MakeNode(absFilePath, Data{
		filePath: absFilePath,
	}), nil
}
