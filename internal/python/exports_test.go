package python

import (
	"context"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"dep-tree/internal/language"
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
					Id:    path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "folder", Alias: "foo_3"}},
					Id:    path.Join(exportsTestFolder, "main.py"),
				},
				//{
				//	Names: []language.ExportName{{Original: "foo"}},
				//	Id:    path.Join(exportsTestFolder, "main.py"),
				// },
				//{
				//	Names: []language.ExportName{{Original: "folder", Alias: "foo"}},
				//	Id:    path.Join(exportsTestFolder, "main.py"),
				// },
				//{
				//	Names: []language.ExportName{{Original: "bar"}},
				//	Id:    path.Join(exportsTestFolder, "foo.py"),
				// },
				//{
				//	Names: []language.ExportName{{Original: "baz", Alias: "baz_2"}},
				//	Id:    path.Join(exportsTestFolder, "folder", "foo.py"),
				// },
				{
					Names: []language.ExportName{{Original: "a"}},
					Id:    path.Join(exportsTestFolder, "main.py"),
				},
				{
					All:   true,
					Names: make([]language.ExportName, 0),
					Id:    path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "module"}},
					Id:    path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "foo"}},
					Id:    path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "foo_1"}, {Original: "foo_2"}},
					Id:    path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "foo_3"}, {Original: "foo_4"}},
					Id:    path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "foo_5"}, {Original: "foo_6"}},
					Id:    path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "func"}},
					Id:    path.Join(exportsTestFolder, "main.py"),
				},
				{
					Names: []language.ExportName{{Original: "Class"}},
					Id:    path.Join(exportsTestFolder, "main.py"),
				},
				//{
				//	Names: []language.ExportName{{Original: "collections", Alias: "collections_abc"}},
				//	Id:    path.Join(exportsTestFolder, "main.py"),
				// },
				//{
				//	Names: []language.ExportName{{Original: "collections", Alias: "collections_abc"}},
				//	Id:    path.Join(exportsTestFolder, "main.py"),
				// },
				{
					Names: []language.ExportName{
						{Original: "a"},
						{Original: "b"},
						{Original: "c"},
					},
					Id: path.Join(exportsTestFolder, "foo.py"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			lang, err := MakePythonLanguage(path.Join(exportsTestFolder, tt.Entrypoint))
			a.NoError(err)

			parsed, err := lang.ParseFile(path.Join(exportsTestFolder, tt.File))
			a.NoError(err)

			_, result, err := lang.ParseExports(context.Background(), parsed)
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
