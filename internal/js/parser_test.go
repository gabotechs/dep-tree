package js

import (
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
	node, err := parser.Entrypoint(id)
	a.NoError(err)
	a.Equal(node.Id, absPath)
	a.Equal(node.Data.dirname, path.Dir(absPath))
	a.Equal(node.Data.content, []byte("console.log(\"hello world!\")\n"))
}

func TestParser_Deps(t *testing.T) {
	tests := []struct {
		Name     string
		Files    map[string]string
		Expected []string
	}{
		{
			Name: "with-imports",
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
			node, err := parser.Entrypoint(id)
			a.NoError(err)
			deps, err := parser.Deps(node)
			a.NoError(err)
			result := make([]string, 0)
			for _, dep := range deps {
				display := parser.Display(dep, node)
				result = append(result, display)
			}

			a.Equal(tt.Expected, result)
		})
	}
}
