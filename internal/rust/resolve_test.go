package rust

import (
	"context"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDirToModChain(t *testing.T) {
	tests := []struct {
		Name     string
		Path     string
		Expected []string
	}{
		{
			Name:     "Simple",
			Path:     path.Join(testFolder, "src", "random", "slice"),
			Expected: []string{"random", "slice"},
		},
		{
			Name:     "Does not output just a .",
			Path:     path.Join(testFolder, "src"),
			Expected: []string{},
		},
		{
			Name:     "src/lib.rs is the source",
			Path:     path.Join(testFolder, "src", "lib.rs"),
			Expected: []string{},
		},
		{
			Name:     "mod file refers to parent folder module",
			Path:     path.Join(testFolder, "src", "div", "mod.rs"),
			Expected: []string{"div"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_, _lang, err := MakeRustLanguage(context.Background(), path.Join(testFolder, "src", "lib.rs"))
			a.NoError(err)

			lang := _lang.(*Language)

			abs, err := filepath.Abs(tt.Path)
			a.NoError(err)

			slices, err := lang.filePathToModChain(abs)
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
			FilePath: path.Join(testFolder, "src", "abs", "abs.rs"),
			Expected: path.Join(testFolder, "src", "abs.rs"),
		},
		{
			Name:     "crate abs abs",
			FilePath: path.Join(testFolder, "src", "div", "div_2.rs"),
			Expected: path.Join(testFolder, "src", "abs", "abs.rs"),
		},
		{
			Name:     "self avg_2 avg",
			FilePath: path.Join(testFolder, "src", "lib.rs"),
			Expected: path.Join(testFolder, "src", "avg_2.rs"),
		},
		{
			Name:     "super div_2 div_2",
			FilePath: path.Join(testFolder, "src", "div", "div.rs"),
			Expected: path.Join(testFolder, "src", "div", "div_2", "div_2.rs"),
		},
		{
			Name:     "self sum",
			FilePath: path.Join(testFolder, "src", "lib.rs"),
			Expected: path.Join(testFolder, "src", "sum.rs"),
		},
		{
			Name:     "crate div div",
			FilePath: path.Join(testFolder, "src", "lib.rs"),
			Expected: path.Join(testFolder, "src", "div", "div.rs"),
		},
		{
			Name:     "un_existing",
			FilePath: path.Join(testFolder, "src", "lib.rs"),
			Expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_, _lang, err := MakeRustLanguage(context.Background(), path.Join(testFolder, "src", "lib.rs"))
			a.NoError(err)

			lang := _lang.(*Language)

			abs, err := filepath.Abs(tt.FilePath)
			a.NoError(err)

			resolved, err := lang.resolve(strings.Split(tt.Name, " "), abs)
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
			FilePath: path.Join(testFolder, "src", "lib.rs"),
			Expected: "could not find mod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_, _lang, err := MakeRustLanguage(context.Background(), path.Join(testFolder, "src", "lib.rs"))
			a.NoError(err)

			lang := _lang.(*Language)

			abs, err := filepath.Abs(tt.FilePath)
			a.NoError(err)

			_, err = lang.resolve(strings.Split(tt.Name, " "), abs)
			a.ErrorContains(err, tt.Expected)
		})
	}
}
