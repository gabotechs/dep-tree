package js

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"dep-tree/internal/language"
)

const importsTestFolder = ".imports_test"

func TestParser_parseImports(t *testing.T) {
	wd, _ := os.Getwd()

	tests := []struct {
		Name           string
		File           string
		Expected       map[string]language.ImportEntry
		ExpectedErrors []string
	}{
		{
			Name: "test 1",
			File: path.Join(importsTestFolder, "index.ts"),
			Expected: map[string]language.ImportEntry{
				path.Join(wd, importsTestFolder, "2", "2.ts"):      {Names: []string{"a", "b"}},
				path.Join(wd, importsTestFolder, "2", "index.ts"):  {All: true},
				path.Join(wd, importsTestFolder, "1", "a", "a.ts"): {All: true},
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

			results, err := lang.ParseImports(parsed)
			a.NoError(err)
			for expectedPath, expectedNames := range tt.Expected {
				resultNames, ok := results.Imports.Get(expectedPath)
				a.Equal(true, ok)
				a.Equal(expectedNames, resultNames)
			}

			a.Equal(len(tt.ExpectedErrors), len(results.Errors))
			if results.Errors != nil {
				for i, err := range results.Errors {
					a.ErrorContains(err, tt.ExpectedErrors[i])
				}
			}
		})
	}
}
