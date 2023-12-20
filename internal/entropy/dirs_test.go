package entropy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_dirs(t *testing.T) {
	tests := []struct {
		Name     string
		Expected []string
	}{
		{
			Name:     "foo/bar/baz",
			Expected: []string{"foo/bar/baz", "foo/bar", "foo"},
		},
		{
			Name:     "foo/bar/baz/",
			Expected: []string{"foo/bar/baz", "foo/bar", "foo"},
		},
		{
			Name:     "/foo/bar/baz/",
			Expected: []string{"foo/bar/baz", "foo/bar", "foo"},
		},
		{
			Name:     "/foo",
			Expected: []string{"foo"},
		},
		{
			Name:     "foo",
			Expected: []string{"foo"},
		},
		{
			Name:     "foo/",
			Expected: []string{"foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			dirTree := NewDirTree()
			actual := dirTree.AddDirs(tt.Name)
			a.Equal(tt.Expected, actual)
		})
	}
}
