package rust

import (
	"context"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/language"
)

func TestLanguage_ParseImports(t *testing.T) {
	absTestFolder, _ := filepath.Abs(path.Join(testFolder))

	tests := []struct {
		Name     string
		Expected []language.ImportEntry
		Errors   []error
	}{
		{
			Name: "lib.rs",
			Expected: []language.ImportEntry{
				{
					All:   true,
					Names: []string{"sum"},
					Path:  path.Join(absTestFolder, "src", "sum.rs"),
				},
				{
					All:   true,
					Names: []string{"div"},
					Path:  path.Join(absTestFolder, "src", "div", "mod.rs"),
				},
				{
					All:   true,
					Names: []string{"avg"},
					Path:  path.Join(absTestFolder, "src", "avg.rs"),
				},
				{
					All:   true,
					Names: []string{"abs"},
					Path:  path.Join(absTestFolder, "src", "abs.rs"),
				},
				{
					All:   true,
					Names: []string{"avg_2"},
					Path:  path.Join(absTestFolder, "src", "avg_2.rs"),
				},
				{
					Names: []string{"abs"},
					Path:  path.Join(absTestFolder, "src", "abs", "abs.rs"),
				},
				{
					Names: []string{"div"},
					Path:  path.Join(absTestFolder, "src", "div", "mod.rs"),
				},
				{
					Names: []string{"avg"},
					Path:  path.Join(absTestFolder, "src", "avg_2.rs"),
				},
				{
					Names: []string{"sum"},
					Path:  path.Join(absTestFolder, "src", "lib.rs"),
				},
				{
					All:  true,
					Path: path.Join(absTestFolder, "src", "sum.rs"),
				},
				{
					Names: []string{"run"},
					Path:  path.Join(absTestFolder, "src", "lib.rs"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_, _lang, err := MakeRustLanguage(context.Background(), path.Join(testFolder, "src", "lib.rs"))
			a.NoError(err)

			lang := _lang.(*Language)

			file, err := lang.ParseFile(path.Join(absTestFolder, "src", tt.Name))
			a.NoError(err)

			exports, err := lang.ParseImports(file)
			a.NoError(err)
			a.Equal(tt.Expected, exports.Imports)
			a.Equal(tt.Errors, exports.Errors)
		})
	}
}
