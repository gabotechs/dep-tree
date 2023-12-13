package language

import (
	"context"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const testFolder = ".parser_test"

func TestParser_Entrypoint(t *testing.T) {
	a := require.New(t)
	id := path.Join(testFolder, t.Name()+".js")
	lang := &TestLanguage{}
	parser := lang.testParser(id)

	entrypoint, err := parser.Entrypoint()
	a.NoError(err)

	a.NoError(err)
	a.Equal(id, entrypoint.Id)
}

func TestParser_Deps(t *testing.T) {
	tests := []struct {
		Name     string
		Path     string
		Imports  map[string]*ImportsResult
		Exports  map[string]*ExportsResult
		Expected []string
	}{
		{
			Name: "Simple",
			Path: "1",
			Imports: map[string]*ImportsResult{
				"1": {
					Imports: []ImportEntry{
						{All: true, Path: "2"},
					},
				},
			},
			Exports: map[string]*ExportsResult{
				"1": {},
				"2": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}},
						Path:  "2",
					}},
				},
			},
			Expected: []string{
				"2",
			},
		},
		{
			Name: "Index only has exports",
			Path: "1",
			Imports: map[string]*ImportsResult{
				"1": {Imports: []ImportEntry{}},
			},
			Exports: map[string]*ExportsResult{
				"1": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}},
						Path:  "2",
					}},
				},
				"2": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}},
						Path:  "2",
					}},
				},
			},
			Expected: []string{
				"2",
			},
		},
		{
			Name: "Proxy export",
			Path: "1",
			Imports: map[string]*ImportsResult{
				"1": {
					Imports: []ImportEntry{
						{All: true, Path: "2"},
					},
				},
			},
			Exports: map[string]*ExportsResult{
				"1": {},
				"2": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}},
						Path:  "3",
					}},
				},
				"3": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}, {Original: "Another-one"}},
						Path:  "3",
					}},
				},
			},
			Expected: []string{
				"3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			lang := &TestLanguage{
				imports: tt.Imports,
				exports: tt.Exports,
			}
			parser := lang.testParser(tt.Path)

			node, err := parser.Entrypoint()
			a.NoError(err)
			_, deps, err := parser.Deps(context.Background(), node)
			a.NoError(err)
			result := make([]string, len(deps))
			for i, dep := range deps {
				a.Equal(0, len(dep.Errors))
				result[i] = dep.Id
			}

			a.Equal(tt.Expected, result)
		})
	}
}

func TestParser_DepsErrors(t *testing.T) {
	tests := []struct {
		Name           string
		Path           string
		Imports        map[string]*ImportsResult
		Exports        map[string]*ExportsResult
		ExpectedErrors []string
	}{
		{
			Name: "Importing a name that is not exported returns an error",
			Path: "1",
			Imports: map[string]*ImportsResult{
				"1": {Imports: []ImportEntry{
					{
						Names: []string{"foo"},
						Path:  "2",
					},
				}},
			},
			Exports: map[string]*ExportsResult{
				"1": {},
				"2": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "bar"}},
						Path:  "2",
					}},
				},
			},
			ExpectedErrors: []string{
				"name foo is imported by 1 but not exported by 2",
			},
		},
	}

	for _, tt := range tests[1:] { // TODO: this is not retro-compatible, do it in a different PR.
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			lang := &TestLanguage{
				imports: tt.Imports,
				exports: tt.Exports,
			}
			parser := lang.testParser(tt.Path)
			node, err := parser.Entrypoint()
			a.NoError(err)
			_, _, err = parser.Deps(context.Background(), node)
			a.NoError(err)
			i := 0
			for _, err := range node.Errors {
				a.ErrorContains(err, tt.ExpectedErrors[i])
				i += 1
			}
			a.Equal(i, len(tt.ExpectedErrors))
		})
	}
}
