package js

import (
	"os"
	"path"

	"dep-tree/internal/graph/node"
)

var Extensions = []string{
	"js", "ts", "tsx", "jsx", "d.ts",
}

type Data struct {
	dirname string
	content []byte
}

func MakeJsNode(absFilePath string) (*node.Node[Data], error) {
	content, err := os.ReadFile(absFilePath)
	if err != nil {
		return nil, err
	}

	dirname := path.Dir(absFilePath)

	return node.MakeNode(absFilePath, Data{
		dirname: dirname,
		content: content,
	}), nil
}
