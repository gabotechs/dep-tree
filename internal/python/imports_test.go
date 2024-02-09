package python

import (
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
		Name                      string
		File                      string
		Entrypoint                string
		Expected                  []language.ImportEntry
		ExpectedErrors            []string
		ExcludeConditionalImports bool
	}{
		{
			Name:       "main.py",
			File:       "main.py",
			Entrypoint: "main.py",
			Expected: []language.ImportEntry{
				language.EmptyImport(filepath.Join(importsTestFolder, "src", "foo.py")),
				language.EmptyImport(filepath.Join(importsTestFolder, "src", "main.py")),
				language.EmptyImport(filepath.Join(importsTestFolder, "src", "main.py")),
				language.EmptyImport(filepath.Join(importsTestFolder, "src", "module", "__init__.py")),
				language.NamesImport([]string{"main"}, filepath.Join(importsTestFolder, "src", "main.py")),
				language.EmptyImport(filepath.Join(importsTestFolder, "src", "main.py")),
				language.NamesImport([]string{"main"}, filepath.Join(importsTestFolder, "src", "main.py")),
				language.AllImport(filepath.Join(importsTestFolder, "src", "module", "__init__.py")),
				language.EmptyImport(filepath.Join(importsTestFolder, "src", "module", "module.py")),
				language.NamesImport([]string{"bar"}, filepath.Join(importsTestFolder, "src", "module", "__init__.py")),
			},
			ExpectedErrors: []string{
				"cannot import file src.py from directory",
				"cannot import file un_existing.py from directory",
			},
		},
		{
			Name:                      "main.py",
			File:                      "main.py",
			Entrypoint:                "main.py",
			ExcludeConditionalImports: true,
			Expected: []language.ImportEntry{
				language.EmptyImport(filepath.Join(importsTestFolder, "src", "foo.py")),
				language.EmptyImport(filepath.Join(importsTestFolder, "src", "main.py")),
				language.EmptyImport(filepath.Join(importsTestFolder, "src", "main.py")),
				language.EmptyImport(filepath.Join(importsTestFolder, "src", "module", "__init__.py")),
				// language.NamesImport([]string{"main"}, filepath.Join(importsTestFolder, "src", "main.py")),
				// language.EmptyImport(filepath.Join(importsTestFolder, "src", "main.py")),
				language.NamesImport([]string{"main"}, filepath.Join(importsTestFolder, "src", "main.py")),
				language.AllImport(filepath.Join(importsTestFolder, "src", "module", "__init__.py")),
				language.EmptyImport(filepath.Join(importsTestFolder, "src", "module", "module.py")),
				language.NamesImport([]string{"bar"}, filepath.Join(importsTestFolder, "src", "module", "__init__.py")),
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
			lang, err := MakePythonLanguage(&Config{
				ExcludeConditionalImports: tt.ExcludeConditionalImports,
			})
			a.NoError(err)

			parsed, err := lang.ParseFile(filepath.Join(importsTestFolder, tt.File))
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
			lang, err := MakePythonLanguage(nil)
			a.NoError(err)

			file := language.FileInfo{Content: &tt.File} //nolint:gosec

			result, err := lang.ParseImports(&file)
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
