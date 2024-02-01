package entropy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDirTree_baseFunctions(t *testing.T) {
	tests := []struct {
		Name              string
		ExpectedFullPaths []string
		ExpectedBaseNames []string
		ExpectedTree      map[string]any
	}{
		{
			Name:              "foo/bar/baz",
			ExpectedFullPaths: []string{"foo/bar/baz", "foo/bar", "foo"},
			ExpectedBaseNames: []string{"foo", "bar", "baz"},
			ExpectedTree: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"baz": map[string]any{},
					},
				},
			},
		},
		{
			Name:              "foo/bar/baz/",
			ExpectedFullPaths: []string{"foo/bar/baz", "foo/bar", "foo"},
			ExpectedBaseNames: []string{"foo", "bar", "baz"},
			ExpectedTree: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"baz": map[string]any{},
					},
				},
			},
		},
		{
			Name:              "/foo/bar/baz/",
			ExpectedFullPaths: []string{"foo/bar/baz", "foo/bar", "foo"},
			ExpectedBaseNames: []string{"foo", "bar", "baz"},
			ExpectedTree: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"baz": map[string]any{},
					},
				},
			},
		},
		{
			Name:              "/foo",
			ExpectedFullPaths: []string{"foo"},
			ExpectedBaseNames: []string{"foo"},
			ExpectedTree: map[string]any{
				"foo": map[string]any{},
			},
		},
		{
			Name:              "foo",
			ExpectedFullPaths: []string{"foo"},
			ExpectedBaseNames: []string{"foo"},
			ExpectedTree: map[string]any{
				"foo": map[string]any{},
			},
		},
		{
			Name:              "foo/",
			ExpectedFullPaths: []string{"foo"},
			ExpectedBaseNames: []string{"foo"},
			ExpectedTree: map[string]any{
				"foo": map[string]any{},
			},
		},
		{
			Name:              "../foo",
			ExpectedFullPaths: []string{"../foo", ".."},
			ExpectedBaseNames: []string{"..", "foo"},
			ExpectedTree: map[string]any{
				"..": map[string]any{
					"foo": map[string]any{},
				},
			},
		},
		{
			Name:              "../../foo",
			ExpectedFullPaths: []string{"../../foo", "../..", ".."},
			ExpectedBaseNames: []string{"..", "..", "foo"},
			ExpectedTree: map[string]any{
				"..": map[string]any{
					"..": map[string]any{
						"foo": map[string]any{},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			dirTree := NewDirTree()
			dirTree.AddDirs(splitBaseNames(tt.Name))
			a.Equal(tt.ExpectedFullPaths, splitFullPaths(tt.Name))
			a.Equal(tt.ExpectedBaseNames, splitBaseNames(tt.Name))
			a.Equal(tt.ExpectedTree, unwrapDirTree(dirTree))
		})
	}
}

func TestDirTree_GroupingsForDir(t *testing.T) {
	tests := []struct {
		Name              string
		Paths             []string
		ExpectedTree      map[string]any
		ExpectedGroupings [][]string
		ExpectedColors    [][]float64
	}{
		{
			Name:  "Single File",
			Paths: []string{"foo/bar"},
			ExpectedTree: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{},
				},
			},
			ExpectedGroupings: [][]string{nil},
			ExpectedColors:    [][]float64{{0, 0, 1}},
		},
		{
			Name:  "two unrelated files",
			Paths: []string{"foo/bar", "baz/bar"},
			ExpectedTree: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{},
				},
				"baz": map[string]any{
					"bar": map[string]any{},
				},
			},
			ExpectedGroupings: [][]string{{"foo"}, {"baz"}},
			ExpectedColors:    [][]float64{{0, 0.76, 1}, {180, 0.76, 1}},
		},
		{
			Name:  "two files with a shared first folder",
			Paths: []string{"foo/bar", "foo/baz"},
			ExpectedTree: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{},
					"baz": map[string]any{},
				},
			},
			ExpectedGroupings: [][]string{{"foo/bar"}, {"foo/baz"}},
			ExpectedColors:    [][]float64{{0, 0.76, 1}, {180, 0.76, 1}},
		},
		{
			Name:  "with middle folders",
			Paths: []string{"foo/bar/baz/1", "foo/bar/baz/2", "bar/foo"},
			ExpectedTree: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"baz": map[string]any{
							"1": map[string]any{},
							"2": map[string]any{},
						},
					},
				},
				"bar": map[string]any{
					"foo": map[string]any{},
				},
			},
			ExpectedGroupings: [][]string{
				{"foo", "foo/bar/baz/1"},
				{"foo", "foo/bar/baz/2"},
				{"bar"},
			},
			ExpectedColors: [][]float64{
				{0, 0.592, 1},
				{180, 0.592, 1},
				{180, 0.76, 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			dirTree := NewDirTree()
			for _, path := range tt.Paths {
				dirTree.AddDirs(splitBaseNames(path))
			}
			a.Equal(tt.ExpectedTree, unwrapDirTree(dirTree))
			var groupings [][]string
			var colors [][]float64
			for _, path := range tt.Paths {
				groupings = append(groupings, dirTree.GroupingsForDir(splitBaseNames(path)))
				colors = append(colors, dirTree.ColorForDir(splitBaseNames(path), HSV))
			}
			a.Equal(tt.ExpectedGroupings, groupings)
			a.Equal(tt.ExpectedColors, colors)
		})
	}
}

func unwrapDirTree(tree *DirTree) interface{} {
	if tree == nil {
		return nil
	}
	result := map[string]any{}
	for el := tree.inner().Front(); el != nil; el = el.Next() {
		result[el.Key] = unwrapDirTree(el.Value.entry)
	}
	return result
}
