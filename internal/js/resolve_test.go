package js

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const resolverTestFolder = ".resolve_test"

func TestFileInfo(t *testing.T) {
	cwd, _ := os.Getwd()
	tests := []struct {
		Name            string
		Entrypoint      string
		ExpectedImports []*Import
		ExpectedExports []*Export
	}{
		{
			Name:       "test 1",
			Entrypoint: path.Join(resolverTestFolder, "test_1", "src", "index.ts"),
			ExpectedImports: []*Import{
				{
					AbsPath: path.Join(cwd, resolverTestFolder, "test_1", "src", "foo.ts"),
				},
				{
					AbsPath: path.Join(cwd, resolverTestFolder, "test_1", "src", "utils", "sum.ts"),
				},
				{
					AbsPath: path.Join(cwd, resolverTestFolder, "test_1", "src", "helpers", "diff.ts"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			parser, err := MakeJsParser(tt.Entrypoint)
			a.NoError(err)
			content, err := os.ReadFile(tt.Entrypoint)
			a.NoError(err)
			dirname := path.Dir(tt.Entrypoint)
			fileInfo, err := parser.ParseFileInfo(content, dirname)
			a.NoError(err)
			expectedImports := make([]string, len(tt.ExpectedImports))
			for i, expectedImport := range tt.ExpectedImports {
				expectedImports[i] = expectedImport.AbsPath
			}

			actualImports := make([]string, len(fileInfo.imports))
			for i, actualImport := range fileInfo.imports {
				actualImports[i] = actualImport.AbsPath
			}

			a.Equal(expectedImports, actualImports)
		})
	}
}
