package js

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testFolder = "parser_test"

func TestParser_Parse(t *testing.T) {
	a := require.New(t)
	id := path.Join(testFolder, t.Name()+".js")
	absPath, err := filepath.Abs(id)
	a.NoError(err)

	node, err := Parser.Parse(id)
	a.NoError(err)
	a.Equal(node.Id, absPath)
	a.Equal(node.Data.dirname, path.Dir(absPath))
	a.Equal(node.Data.content, []byte("console.log(\"hello world!\")\n"))
}

func TestParser_Deps(t *testing.T) {
	thisDir, _ := os.Getwd()

	tests := []struct {
		Name      string
		Expected  []string
		Normalize bool
	}{
		{
			Name: "deps",
			Expected: []string{
				"geometries/Geometries.js",
				"geometries/Geometries.js",
				"parser_test/.export",
				"parser_test/views/ListView",
				"parser_test/views/AddView",
				"parser_test/views/EditView",
				"parser_test/views/ListView",
				"parser_test/views/AddView",
				"parser_test/views/EditView",
			},
		},
		{
			Name:      "with-imports",
			Normalize: true,
			Expected: []string{
				"parser_test/with-imports-imported/imported.js",
			},
		},
		{
			Name:      "with-imports-index",
			Normalize: true,
			Expected: []string{
				"parser_test/with-imports-index-imported/other.js",
				"parser_test/with-imports-index-imported/one.js",
				"parser_test/with-imports-index-imported/index.js",
			},
		},

		{
			Name: "custom-1",
			Expected: []string{
				"@parsers/DateSchema",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			id := path.Join(testFolder, path.Base(t.Name())+".js")

			node, err := Parser.Parse(id)
			a.NoError(err)
			deps := Parser.Deps(node)
			result := make([]string, 0)
			for _, dep := range deps {
				if tt.Normalize {
					dep, err = normalizeId(dep)
					a.NoError(err)
				}
				result = append(result, strings.ReplaceAll(dep, thisDir+"/", ""))
			}

			a.Equal(tt.Expected, result)
		})
	}
}
