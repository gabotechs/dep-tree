package rust

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"dep-tree/internal/language"
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
					Id:    path.Join(absTestFolder, "src", "sum.rs"),
				},
				{
					All:   true,
					Names: []string{"div"},
					Id:    path.Join(absTestFolder, "src", "div", "mod.rs"),
				},
				{
					All:   true,
					Names: []string{"avg"},
					Id:    path.Join(absTestFolder, "src", "avg.rs"),
				},
				{
					All:   true,
					Names: []string{"abs"},
					Id:    path.Join(absTestFolder, "src", "abs.rs"),
				},
				{
					All:   true,
					Names: []string{"avg_2"},
					Id:    path.Join(absTestFolder, "src", "avg_2.rs"),
				},
				{
					Names: []string{"abs"},
					Id:    path.Join(absTestFolder, "src", "abs", "abs.rs"),
				},
				{
					Names: []string{"div"},
					Id:    path.Join(absTestFolder, "src", "div", "mod.rs"),
				},
				{
					Names: []string{"avg"},
					Id:    path.Join(absTestFolder, "src", "avg_2.rs"),
				},
				{
					Names: []string{"sum"},
					Id:    path.Join(absTestFolder, "src", "lib.rs"),
				},
				{
					All:   true,
					Names: []string{},
					Id:    path.Join(absTestFolder, "src", "sum.rs"),
				},
				{
					Names: []string{"run"},
					Id:    path.Join(absTestFolder, "src", "lib.rs"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_lang, err := MakeRustLanguage(path.Join(testFolder, "src", "lib.rs"))
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
