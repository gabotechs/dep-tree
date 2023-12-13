package python

import (
	"context"
	"path"
	"path/filepath"
	"testing"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python/python_grammar"

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
				// {
				//	All: true,
				//	Path:  path.Join(importsTestFolder, "src", "foo.py"),
				// },
				// {
				//	All: true,
				//	Path:  path.Join(importsTestFolder, "src", "main.py"),
				// },
				// {
				//	All: true,
				//	Path:  path.Join(importsTestFolder, "src", "main.py"),
				// },
				// {
				//	All: true,
				//	Path:  path.Join(importsTestFolder, "src", "module", "__init__.py"),
				// },
				// {
				//	Names: []string{"main"},
				//	Path:    path.Join(importsTestFolder, "src", "main.py"),
				// },
				{
					All:  true,
					Path: path.Join(importsTestFolder, "src", "main.py"),
				},
				{
					Names: []string{"main"},
					Path:  path.Join(importsTestFolder, "src", "main.py"),
				},
				{
					All:  true,
					Path: path.Join(importsTestFolder, "src", "module", "__init__.py"),
				},
				{
					All:  true,
					Path: path.Join(importsTestFolder, "src", "module", "module.py"),
				},
				{
					Names: []string{"bar"},
					Path:  path.Join(importsTestFolder, "src", "module", "__init__.py"),
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
			_, lang, err := MakePythonLanguage(context.Background(), path.Join(importsTestFolder, tt.Entrypoint), nil)
			a.NoError(err)

			parsed, err := lang.ParseFile(path.Join(importsTestFolder, tt.File))
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

func TestLanguage_ParseImports_Errors(t *testing.T) {
	tests := []struct {
		Name           string
		File           python_grammar.File
		ExpectedErrors []string
	}{
		{
			Name: "Import Errors",
			File: python_grammar.File{
				Statements: []*python_grammar.Statement{
					nil,
					{
						FromImport: &python_grammar.FromImport{
							Relative: make([]bool, 3),
							Path:     []string{"non-existent"},
						},
					},
				},
			},
			ExpectedErrors: []string{
				"could not resolve relative import",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_, lang, err := MakePythonLanguage(context.Background(), path.Join(importsTestFolder, "main.py"), nil)
			a.NoError(err)

			result, err := lang.ParseImports(&tt.File) //nolint:gosec
			a.NoError(err)

			a.Equal(len(tt.ExpectedErrors), len(result.Errors))
			if result.Errors != nil {
				for i, err := range result.Errors {
					a.ErrorContains(err, tt.ExpectedErrors[i])
				}
			}
		})
	}
}
