package rust

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/language"
)

func TestLanguage_ParseImports(t *testing.T) {
	absTestFolder, _ := filepath.Abs(testFolder)

	tests := []struct {
		Name     string
		Expected []language.ImportEntry
		Errors   []error
	}{
		{
			Name: "lib.rs",
			Expected: []language.ImportEntry{
				{
					All:     true,
					Symbols: []string{"sum"},
					AbsPath: filepath.Join(absTestFolder, "src", "sum.rs"),
				},
				{
					All:     true,
					Symbols: []string{"div"},
					AbsPath: filepath.Join(absTestFolder, "src", "div", "mod.rs"),
				},
				{
					All:     true,
					Symbols: []string{"avg"},
					AbsPath: filepath.Join(absTestFolder, "src", "avg.rs"),
				},
				{
					All:     true,
					Symbols: []string{"abs"},
					AbsPath: filepath.Join(absTestFolder, "src", "abs.rs"),
				},
				{
					All:     true,
					Symbols: []string{"avg_2"},
					AbsPath: filepath.Join(absTestFolder, "src", "avg_2.rs"),
				},
				{
					Symbols: []string{"abs"},
					AbsPath: filepath.Join(absTestFolder, "src", "abs", "abs.rs"),
				},
				{
					Symbols: []string{"div"},
					AbsPath: filepath.Join(absTestFolder, "src", "div", "mod.rs"),
				},
				{
					Symbols: []string{"avg"},
					AbsPath: filepath.Join(absTestFolder, "src", "avg_2.rs"),
				},
				{
					Symbols: []string{"sum"},
					AbsPath: filepath.Join(absTestFolder, "src", "lib.rs"),
				},
				{
					All:     true,
					AbsPath: filepath.Join(absTestFolder, "src", "sum.rs"),
				},
				{
					Symbols: []string{"run"},
					AbsPath: filepath.Join(absTestFolder, "src", "lib.rs"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_lang, err := MakeRustLanguage(nil)
			a.NoError(err)

			lang := _lang.(*Language)

			file, err := lang.ParseFile(filepath.Join(absTestFolder, "src", tt.Name))
			a.NoError(err)

			exports, err := lang.ParseImports(file)
			a.NoError(err)
			a.Equal(tt.Expected, exports.Imports)
			a.Equal(tt.Errors, exports.Errors)
		})
	}
}
