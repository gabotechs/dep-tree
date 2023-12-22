package entropy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_dirs(t *testing.T) {
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
			dirTree.AddDirs(tt.Name)
			a.Equal(tt.ExpectedFullPaths, splitFullPaths(tt.Name))
			a.Equal(tt.ExpectedBaseNames, splitBaseNames(tt.Name))
			a.Equal(tt.ExpectedTree, unwrapDirTree(dirTree))
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
