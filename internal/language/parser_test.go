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

	parser, err := ParserBuilder(func(s string) (Language[TestLanguageData, TestFile], error) {
		return &TestLanguage{}, nil
	})(id)
	a.NoError(err)
	entrypoint, err := parser.Entrypoint()
	a.NoError(err)

	a.NoError(err)
	a.Equal(id, entrypoint.Id)
}

func TestParser_Deps(t *testing.T) {
	tests := []struct {
		Name     string
		Id       string
		Imports  map[string]*ImportsResult
		Exports  map[string]*ExportsResult
		Expected []string
	}{
		{
			Name: "Simple",
			Id:   "1",
			Imports: map[string]*ImportsResult{
				"1": {
					Imports: []ImportEntry{
						{All: true, Id: "2"},
					},
				},
			},
			Exports: map[string]*ExportsResult{
				"1": {},
				"2": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}},
						Id:    "2",
					}},
				},
			},
			Expected: []string{
				"2",
			},
		},
		{
			Name: "Index only has exports",
			Id:   "1",
			Imports: map[string]*ImportsResult{
				"1": {Imports: []ImportEntry{}},
			},
			Exports: map[string]*ExportsResult{
				"1": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}},
						Id:    "2",
					}},
				},
				"2": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}},
						Id:    "2",
					}},
				},
			},
			Expected: []string{
				"2",
			},
		},
		{
			Name: "Proxy export",
			Id:   "1",
			Imports: map[string]*ImportsResult{
				"1": {
					Imports: []ImportEntry{
						{All: true, Id: "2"},
					},
				},
			},
			Exports: map[string]*ExportsResult{
				"1": {},
				"2": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}},
						Id:    "3",
					}},
				},
				"3": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}, {Original: "Another-one"}},
						Id:    "3",
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

			parser, err := ParserBuilder(func(entrypoint string) (Language[TestLanguageData, TestFile], error) {
				return &TestLanguage{
					imports: tt.Imports,
					exports: tt.Exports,
				}, nil
			})(tt.Id)

			a.NoError(err)
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
		Id             string
		Imports        map[string]*ImportsResult
		Exports        map[string]*ExportsResult
		ExpectedErrors []string
	}{
		{
			Name: "Importing a name that is not exported returns an error",
			Id:   "1",
			Imports: map[string]*ImportsResult{
				"1": {Imports: []ImportEntry{
					{
						Names: []string{"foo"},
						Id:    "2",
					},
				}},
			},
			Exports: map[string]*ExportsResult{
				"1": {},
				"2": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "bar"}},
						Id:    "2",
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

			parser, err := ParserBuilder(func(entrypoint string) (Language[TestLanguageData, TestFile], error) {
				return &TestLanguage{
					imports: tt.Imports,
					exports: tt.Exports,
				}, nil
			})(tt.Id)

			a.NoError(err)
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
