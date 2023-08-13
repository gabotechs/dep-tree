package python

import (
	"context"
	"path"
	"path/filepath"
	"testing"

	"dep-tree/internal/language"

	"github.com/stretchr/testify/require"
)

const importsTestFolder = ".imports_test"

func TestLanguage_ParseImports(t *testing.T) {
	importsTestFolder, _ := filepath.Abs(importsTestFolder)

	tests := []struct {
		Name           string
		File           string
		Entrypoint     string
		Expected       []language.ImportEntry
		ExpectedErrors []string
	}{
		{
			Name:       "main.py",
			File:       "main.py",
			Entrypoint: "main.py",
			Expected: []language.ImportEntry{
				{
					Names: []string{"main"},
					Id:    path.Join(importsTestFolder, "src", "main.py"),
				},
				{
					All: true,
					Id:  path.Join(importsTestFolder, "src", "main.py"),
				},
				{
					Names: []string{"main"},
					Id:    path.Join(importsTestFolder, "src", "main.py"),
				},
				{
					All: true,
					Id:  path.Join(importsTestFolder, "src", "module", "__init__.py"),
				},
				{
					All: true,
					Id:  path.Join(importsTestFolder, "src", "module", "module.py"),
				},
			},
			ExpectedErrors: []string{
				"cannot import file src.py from directory",
				"cannot import file un_existing.py from directory",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			lang, err := MakePythonLanguage(path.Join(importsTestFolder, tt.Entrypoint))
			a.NoError(err)

			parsed, err := lang.ParseFile(path.Join(importsTestFolder, tt.File))
			a.NoError(err)

			_, result, err := lang.ParseImports(context.Background(), parsed)
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
