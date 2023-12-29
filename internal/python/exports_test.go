package python

import (
	"context"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/language"
)

const exportsTestFolder = ".exports_test"

func TestLanguage_ParseExports(t *testing.T) {
	exportsTestFolder, _ := filepath.Abs(exportsTestFolder)

	tests := []struct {
		Name           string
		File           string
		Entrypoint     string
		Expected       []language.ExportEntry
		ExpectedErrors []string
	}{
		{
			Name:       "main.py",
			File:       "main.py",
			Entrypoint: "main.py",
			Expected: []language.ExportEntry{
				{
					Names: []language.ExportName{{Original: "foo", Alias: "foo_2"}},
					Path:  path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "foo", Alias: "foo_3"}},
					Path:  path.Join(exportsTestFolder, "main.py"),
				},
				//{
				//	Names: []language.ExportName{{Original: "foo"}},
				//	Path:    path.Join(exportsTestFolder, "main.py"),
				// },
				//{
				//	Names: []language.ExportName{{Original: "folder", Alias: "foo"}},
				//	Path:    path.Join(exportsTestFolder, "main.py"),
				// },
				//{
				//	Names: []language.ExportName{{Original: "bar"}},
				//	Path:    path.Join(exportsTestFolder, "foo.py"),
				// },
				//{
				//	Names: []language.ExportName{{Original: "baz", Alias: "baz_2"}},
				//	Path:    path.Join(exportsTestFolder, "folder", "foo.py"),
				// },
				{
					Names: []language.ExportName{{Original: "a"}},
					Path:  path.Join(exportsTestFolder, "main.py"),
				},
				{
					All:  true,
					Path: path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "module"}},
					Path:  path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "foo"}},
					Path:  path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "foo_1"}, {Original: "foo_2"}},
					Path:  path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "foo_3"}, {Original: "foo_4"}},
					Path:  path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "foo_5"}, {Original: "foo_6"}},
					Path:  path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "func"}},
					Path:  path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "Class"}},
					Path:  path.Join(exportsTestFolder, "main.py"),
				},
				//{
				//	Names: []language.ExportName{{Original: "collections", Alias: "collections_abc"}},
				//	Path:    path.Join(exportsTestFolder, "main.py"),
				// },
				//{
				//	Names: []language.ExportName{{Original: "collections", Alias: "collections_abc"}},
				//	Path:    path.Join(exportsTestFolder, "main.py"),
				// },
				{
					Names: []language.ExportName{
						{Original: "a"},
						{Original: "b"},
						{Original: "c"},
					},
					Path: path.Join(exportsTestFolder, "foo.py"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_, lang, err := MakePythonLanguage(context.Background(), path.Join(exportsTestFolder, tt.Entrypoint), nil)
			a.NoError(err)

			parsed, err := lang.ParseFile(path.Join(exportsTestFolder, tt.File))
			a.NoError(err)

			result, err := lang.ParseExports(parsed)
			a.NoError(err)
			a.Equal(tt.Expected, result.Exports)

			a.Equal(len(tt.ExpectedErrors), len(result.Errors))
			if result.Errors != nil {
				for i, err := range result.Errors {
					a.ErrorContains(err, tt.ExpectedErrors[i])
				}
			}
		})
	}
}
