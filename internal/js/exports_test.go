package js

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/language"
)

const exportsTestFolder = ".exports_test"

func TestParser_parseExports(t *testing.T) {
	cwd, _ := os.Getwd()

	tests := []struct {
		Name           string
		File           string
		Expected       []language.ExportEntry
		ExpectedErrors []string
	}{
		{
			Name: "test",
			File: path.Join(exportsTestFolder, "src", "index.js"),
			Expected: []language.ExportEntry{
				{
					All: true,
					Id:  path.Join(cwd, exportsTestFolder, "src", "utils", "index.js"),
				},
				{
					Names: []language.ExportName{{Original: "Unexisting"}, {Original: "UnSorter", Alias: "UnSorterAlias"}},
					Id:    path.Join(cwd, exportsTestFolder, "src", "utils", "index.js"),
				},
				{
					Names: []language.ExportName{{Original: "aliased"}},
					Id:    path.Join(cwd, exportsTestFolder, "src", "utils", "unsort.js"),
				},
			},
			ExpectedErrors: []string{
				"could not perform relative import for './unexisting'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			lang, err := MakeJsLanguage(tt.File)
			a.NoError(err)

			parsed, err := lang.ParseFile(tt.File)
			a.NoError(err)

			_, exports, err := lang.ParseExports(context.Background(), parsed)
			a.NoError(err)
			a.Equal(tt.Expected, exports.Exports)

			a.Equal(len(tt.ExpectedErrors), len(exports.Errors))
			if exports.Exports != nil {
				for i, err := range exports.Errors {
					a.ErrorContains(err, tt.ExpectedErrors[i])
				}
			}
		})
	}
}
