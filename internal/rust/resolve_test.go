package rust

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testWorkspacesFolder = ".workspaces_project"

func TestDirToModChain(t *testing.T) {
	tests := []struct {
		Name     string
		Path     string
		Expected []string
	}{
		{
			Name:     "Simple",
			Path:     filepath.Join(testFolder, "src", "random", "slice"),
			Expected: []string{"random", "slice"},
		},
		{
			Name:     "Does not output just a .",
			Path:     filepath.Join(testFolder, "src"),
			Expected: []string{},
		},
		{
			Name:     "src/lib.rs is the source",
			Path:     filepath.Join(testFolder, "src", "lib.rs"),
			Expected: []string{},
		},
		{
			Name:     "mod file refers to parent folder module",
			Path:     filepath.Join(testFolder, "src", "div", "mod.rs"),
			Expected: []string{"div"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			entrypoint := filepath.Join(testFolder, "src", "lib.rs")
			a := require.New(t)
			abs, err := filepath.Abs(tt.Path)
			a.NoError(err)

			slices, err := filePathToModChain(abs, entrypoint)
			a.NoError(err)

			a.Equal(tt.Expected, slices)
		})
	}
}

func TestResolve(t *testing.T) {
	tests := []struct {
		Name     string
		FilePath string
		Expected string
	}{
		{
			Name:     "crate abs",
			FilePath: filepath.Join(testFolder, "src", "abs", "abs.rs"),
			Expected: filepath.Join(testFolder, "src", "abs.rs"),
		},
		{
			Name:     "crate abs abs",
			FilePath: filepath.Join(testFolder, "src", "div", "div_2.rs"),
			Expected: filepath.Join(testFolder, "src", "abs", "abs.rs"),
		},
		{
			Name:     "self avg_2 avg",
			FilePath: filepath.Join(testFolder, "src", "lib.rs"),
			Expected: filepath.Join(testFolder, "src", "avg_2.rs"),
		},
		{
			Name:     "super div_2 div_2",
			FilePath: filepath.Join(testFolder, "src", "div", "div.rs"),
			Expected: filepath.Join(testFolder, "src", "div", "div_2", "div_2.rs"),
		},
		{
			Name:     "self sum",
			FilePath: filepath.Join(testFolder, "src", "lib.rs"),
			Expected: filepath.Join(testFolder, "src", "sum.rs"),
		},
		{
			Name:     "crate div div",
			FilePath: filepath.Join(testFolder, "src", "lib.rs"),
			Expected: filepath.Join(testFolder, "src", "div", "div.rs"),
		},
		{
			Name:     "un_existing",
			FilePath: filepath.Join(testFolder, "src", "lib.rs"),
			Expected: "",
		},
		{
			Name:     "foo foo",
			FilePath: filepath.Join(testWorkspacesFolder, "src", "main.rs"),
			Expected: filepath.Join(testWorkspacesFolder, "foo", "src", "lib.rs"),
		},
		{
			Name:     "bar bar",
			FilePath: filepath.Join(testWorkspacesFolder, "src", "main.rs"),
			Expected: filepath.Join(testWorkspacesFolder, "bar", "src", "lib.rs"),
		},
		{
			Name:     "bar",
			FilePath: filepath.Join(testWorkspacesFolder, "foo", "src", "lib.rs"),
			Expected: filepath.Join(testWorkspacesFolder, "bar", "src", "lib.rs"),
		},
		{
			Name:     "baz",
			FilePath: filepath.Join(testWorkspacesFolder, "foo", "src", "lib.rs"),
			Expected: filepath.Join(testWorkspacesFolder, "baz", "src", "lib.rs"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			abs, err := filepath.Abs(tt.FilePath)
			a.NoError(err)

			resolved, err := resolve(strings.Split(tt.Name, " "), abs)
			a.NoError(err)

			var expectedAbs string
			if tt.Expected != "" {
				expectedAbs, err = filepath.Abs(tt.Expected)
				a.NoError(err)
			}

			a.Equal(expectedAbs, resolved)
		})
	}
}

func TestResolveErrors(t *testing.T) {
	tests := []struct {
		Name     string
		FilePath string
		Expected string
	}{
		{
			Name:     "crate un_existing",
			FilePath: filepath.Join(testFolder, "src", "lib.rs"),
			Expected: "could not find mod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			abs, err := filepath.Abs(tt.FilePath)
			a.NoError(err)

			_, err = resolve(strings.Split(tt.Name, " "), abs)
			a.ErrorContains(err, tt.Expected)
		})
	}
}
