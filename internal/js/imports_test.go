package js

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/language"
)

const importsTestFolder = ".imports_test"

func TestParser_parseImports(t *testing.T) {
	wd, _ := os.Getwd()

	tests := []struct {
		Name           string
		File           string
		Expected       []language.ImportEntry
		ExpectedErrors []string
	}{
		{
			Name: "test 1",
			File: filepath.Join(importsTestFolder, "index.ts"),
			Expected: []language.ImportEntry{
				{Names: []string{"a", "b"}, Path: filepath.Join(wd, importsTestFolder, "2", "2.ts")},
				{All: true, Path: filepath.Join(wd, importsTestFolder, "2", "index.ts")},
				{All: true, Path: filepath.Join(wd, importsTestFolder, "1", "a", "a.ts")},
				{All: true, Path: filepath.Join(wd, importsTestFolder, "1", "a", "index.ts")},
				{Names: []string{"Unexisting"}, Path: filepath.Join(wd, importsTestFolder, "1", "a", "index.ts")},
				{All: true, Path: filepath.Join(wd, importsTestFolder, "2", "2.ts")},
				{Names: []string{"a", "b"}, Path: filepath.Join(wd, importsTestFolder, "2", "2.ts")},
				{Path: filepath.Join(wd, importsTestFolder, "1", "a", "index.ts")},
			},
			ExpectedErrors: []string{
				"could not perform relative import for './unexisting'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			lang, err := MakeJsLanguage(tt.File, nil)
			a.NoError(err)

			parsed, err := lang.ParseFile(tt.File)
			a.NoError(err)

			result, err := lang.ParseImports(parsed)
			a.NoError(err)
			a.Equal(tt.Expected, result.Imports)

			a.Equal(len(tt.ExpectedErrors), len(result.Errors))
			if result.Errors != nil {
				for i, err := range result.Errors {
					a.ErrorContains(err, tt.ExpectedErrors[i])
				}
			}
		})
	}
}
