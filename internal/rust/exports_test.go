package rust

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/language"
)

func TestLanguage_ParseExports(t *testing.T) {
	absTestFolder, _ := filepath.Abs(testFolder)

	tests := []struct {
		Name     string
		Expected []language.ExportEntry
		Errors   []error
	}{
		{
			Name: "lib.rs",
			Expected: []language.ExportEntry{
				{
					Symbols: []language.ExportSymbol{{Original: "div"}},
					AbsPath: filepath.Join(absTestFolder, "src", "lib.rs"),
				},
				{
					Symbols: []language.ExportSymbol{{Original: "abs"}},
					AbsPath: filepath.Join(absTestFolder, "src", "abs", "abs.rs"),
				},
				{
					Symbols: []language.ExportSymbol{{Original: "div"}},
					AbsPath: filepath.Join(absTestFolder, "src", "div", "mod.rs"),
				},
				{
					Symbols: []language.ExportSymbol{{Original: "avg"}},
					AbsPath: filepath.Join(absTestFolder, "src", "avg_2.rs"),
				},
				{
					Symbols: []language.ExportSymbol{{Original: "sum"}},
					AbsPath: filepath.Join(absTestFolder, "src", "lib.rs"),
				},
				{
					All:     true,
					AbsPath: filepath.Join(absTestFolder, "src", "sum.rs"),
				},
				{
					Symbols: []language.ExportSymbol{{Original: "run"}},
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

			exports, err := lang.ParseExports(file)
			a.NoError(err)
			a.Equal(tt.Expected, exports.Exports)
			a.Equal(tt.Errors, exports.Errors)
		})
	}
}
