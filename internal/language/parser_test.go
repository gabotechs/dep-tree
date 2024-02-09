package language

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser_shouldExclude(t *testing.T) {
	tests := []struct {
		Name     string
		Paths    []string
		Exclude  []string
		Expected []string
	}{
		{
			Name:     "simple",
			Paths:    []string{"/foo/bar.ts", "/foo/baz.ts"},
			Exclude:  []string{"/foo/bar.ts"},
			Expected: []string{"/foo/baz.ts"},
		},
		{
			Name:     "globstar",
			Paths:    []string{"/foo/bar.ts", "/foo/baz.ts"},
			Exclude:  []string{"/foo/*.ts"},
			Expected: nil,
		},
		{
			Name:     "globstar 2",
			Paths:    []string{"/foo/1/2/3/4/bar.ts", "/foo/baz.ts"},
			Exclude:  []string{"/foo/**/*.ts"},
			Expected: nil,
		},
		{
			Name:     "globstar 3",
			Paths:    []string{"/foo/1/2/3/4/bar.ts", "/foo/baz.ts"},
			Exclude:  []string{"2/**/*.ts"},
			Expected: []string{"/foo/1/2/3/4/bar.ts", "/foo/baz.ts"},
		},
		{
			Name:     "globstar 4",
			Paths:    []string{"/foo/1/2/3/4/bar.ts", "/foo/baz.ts"},
			Exclude:  []string{"*/2/**/*.ts"},
			Expected: []string{"/foo/1/2/3/4/bar.ts", "/foo/baz.ts"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			parser := Parser{Exclude: tt.Exclude}
			var result []string
			for _, path := range tt.Paths {
				if !parser.shouldExclude(path) {
					result = append(result, path)
				}
			}
			a.Equal(tt.Expected, result)
		})
	}
}

func TestParser_Deps(t *testing.T) {
	tests := []struct {
		Name              string
		Path              string
		Imports           map[string]*ImportsResult
		Exports           map[string]*ExportsEntries
		ExpectedUnwrapped []string
		ExpectedWrapped   []string
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
			Exports: map[string]*ExportsEntries{
				"1": {},
				"2": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}},
						Path:  "2",
					}},
				},
			},
			ExpectedUnwrapped: []string{
				"2",
			},
			ExpectedWrapped: []string{
				"2",
			},
		},
		{
			Name: "Index only has exports",
			Path: "1",
			Imports: map[string]*ImportsResult{
				"1": {Imports: []ImportEntry{}},
			},
			Exports: map[string]*ExportsEntries{
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
			ExpectedUnwrapped: []string{
				"2",
			},
			ExpectedWrapped: []string{
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
			Exports: map[string]*ExportsEntries{
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
			ExpectedUnwrapped: []string{
				"3",
			},
			ExpectedWrapped: []string{
				"2",
			},
		},
		{
			Name: "Exports are treated as imports in node entry",
			Path: "1",
			Imports: map[string]*ImportsResult{
				"1": {
					Imports: []ImportEntry{},
				},
			},
			Exports: map[string]*ExportsEntries{
				"1": {
					Exports: []ExportEntry{{
						All:  true,
						Path: "2",
					}, {
						Names: []ExportName{{Original: "Exported-3"}},
						Path:  "3",
					}},
				},
				"2": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}},
						Path:  "4",
					}, {
						Names: []ExportName{{Original: "Exported-2"}},
						Path:  "2",
					}},
				},
				"3": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Another-one", Alias: "Exported-3"}},
						Path:  "4",
					}},
				},
				"4": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}, {Original: "Another-one"}},
						Path:  "4",
					}},
				},
			},
			ExpectedUnwrapped: []string{
				"4", "2",
			},
			ExpectedWrapped: []string{
				"2", "3",
			},
		},
		{
			Name: "Exports are treated as imports in node entry (2)",
			Path: "1",
			Imports: map[string]*ImportsResult{
				"1": {
					Imports: []ImportEntry{},
				},
			},
			Exports: map[string]*ExportsEntries{
				"1": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported-3"}},
						Path:  "3",
					}},
				},
				"2": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}},
						Path:  "4",
					}, {
						Names: []ExportName{{Original: "Exported-2"}},
						Path:  "2",
					}},
				},
				"3": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Another-one", Alias: "Exported-3"}},
						Path:  "4",
					}},
				},
				"4": {
					Exports: []ExportEntry{{
						Names: []ExportName{{Original: "Exported"}, {Original: "Another-one"}},
						Path:  "4",
					}},
				},
			},
			ExpectedUnwrapped: []string{
				"4",
			},
			ExpectedWrapped: []string{
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
			parser := lang.testParser()
			parser.UnwrapProxyExports = true
			node, err := parser.Node(tt.Path)
			a.NoError(err)
			deps, err := parser.Deps(node)
			a.NoError(err)
			result := make([]string, len(deps))
			for i, dep := range deps {
				a.Equal(0, len(dep.Errors))
				result[i] = dep.Id
			}
			a.Equal(tt.ExpectedUnwrapped, result)

			parser.UnwrapProxyExports = false

			deps, err = parser.Deps(node)
			a.NoError(err)
			result = make([]string, len(deps))
			for i, dep := range deps {
				a.Equal(0, len(dep.Errors))
				result[i] = dep.Id
			}
			a.Equal(tt.ExpectedWrapped, result)
		})
	}
}

func TestParser_DepsErrors(t *testing.T) {
	tests := []struct {
		Name           string
		Path           string
		Imports        map[string]*ImportsResult
		Exports        map[string]*ExportsEntries
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
			Exports: map[string]*ExportsEntries{
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
			parser := lang.testParser()
			node, err := parser.Node(tt.Path)
			a.NoError(err)
			_, err = parser.Deps(node)
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
