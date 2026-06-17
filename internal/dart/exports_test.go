package dart

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/language"
)

func TestLanguage_ParseExports(t *testing.T) {
	absTestFolder, _ := filepath.Abs(testFolder)
	libDir := filepath.Join(absTestFolder, "lib")

	tests := []struct {
		Name     string
		File     string
		Expected []language.ExportEntry
		Errors   []error
	}{
		{
			Name: "main.dart",
			File: filepath.Join("lib", "main.dart"),
			Expected: []language.ExportEntry{
				// export 'src/bar.dart' show Bar;
				{
					Symbols: []language.ExportSymbol{{Original: "Bar"}},
					AbsPath: filepath.Join(libDir, "src", "bar.dart"),
				},
			},
		},
		{
			Name: "src/bar.dart",
			File: filepath.Join("lib", "src", "bar.dart"),
			Expected: []language.ExportEntry{
				{
					Symbols: []language.ExportSymbol{{Original: "Bar"}},
					AbsPath: filepath.Join(libDir, "src", "bar.dart"),
				},
				{
					Symbols: []language.ExportSymbol{{Original: "Status"}},
					AbsPath: filepath.Join(libDir, "src", "bar.dart"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_lang, err := MakeDartLanguage(nil)
			a.NoError(err)
			lang := _lang.(*Language)

			file, err := lang.ParseFile(filepath.Join(absTestFolder, tt.File))
			a.NoError(err)

			exports, err := lang.ParseExports(file)
			a.NoError(err)
			a.Equal(tt.Expected, exports.Exports)
			a.Equal(tt.Errors, exports.Errors)
		})
	}
}
