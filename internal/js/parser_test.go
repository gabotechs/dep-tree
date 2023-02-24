package js

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testFolder = ".parser_test"

func TestParser_Entrypoint(t *testing.T) {
	a := require.New(t)
	id := path.Join(testFolder, t.Name()+".js")
	absPath, err := filepath.Abs(id)
	a.NoError(err)

	parser, err := MakeJsParser(id)
	a.NoError(err)
	node, err := parser.Entrypoint()
	a.NoError(err)
	a.Equal(absPath, node.Id)
}

func TestParser_Deps(t *testing.T) {
	tests := []struct {
		Name     string
		Expected []string
	}{
		{
			Name: "with-imports",
			Expected: []string{
				"with-imports-imported/imported.js",
			},
		},
		{
			Name: "with-exports",
			Expected: []string{
				"with-imports-imported/imported.js",
			},
		},
		{
			Name: "with-imports-index",
			Expected: []string{
				"with-imports-index-imported/other.js",
				"with-imports-index-imported/one.js",
				"with-imports-index-imported/index.js",
			},
		},
		{
			Name: "with-imports-nested",
			Expected: []string{
				"generated/generated.js",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			id := path.Join(testFolder, path.Base(t.Name())+".js")
			if _, err := os.Stat(id); err != nil {
				id = path.Join(testFolder, path.Base(t.Name()), "index.js")
			}

			parser, err := MakeJsParser(id)
			a.NoError(err)
			node, err := parser.Entrypoint()
			a.NoError(err)
			_, deps, err := parser.Deps(context.Background(), node)
			a.NoError(err)
			result := make([]string, len(deps))
			for i, dep := range deps {
				result[i] = parser.Display(dep)
			}

			a.Equal(tt.Expected, result)
		})
	}
}
