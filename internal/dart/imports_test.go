package dart

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/language"
)

const testFolder = ".sample_project"

func TestLanguage_ParseImports(t *testing.T) {
	absTestFolder, _ := filepath.Abs(testFolder)
	libDir := filepath.Join(absTestFolder, "lib")

	tests := []struct {
		Name     string
		File     string
		Expected []language.ImportEntry
		Errors   []error
	}{
		{
			Name: "main.dart",
			File: filepath.Join("lib", "main.dart"),
			Expected: []language.ImportEntry{
				// import 'package:sample_app/foo.dart';
				{
					All:     true,
					AbsPath: filepath.Join(libDir, "foo.dart"),
				},
				// import 'utils/helper.dart' show helper;
				{
					Symbols: []string{"helper"},
					AbsPath: filepath.Join(libDir, "utils", "helper.dart"),
				},
				// import 'src/bar.dart' as bar;
				{
					All:     true,
					AbsPath: filepath.Join(libDir, "src", "bar.dart"),
				},
				// export 'src/bar.dart' show Bar;
				{
					Symbols: []string{"Bar"},
					AbsPath: filepath.Join(libDir, "src", "bar.dart"),
				},
				// part 'main.g.dart';
				{
					All:     true,
					AbsPath: filepath.Join(libDir, "main.g.dart"),
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

			imports, err := lang.ParseImports(file)
			a.NoError(err)
			a.Equal(tt.Expected, imports.Imports)
			a.Equal(tt.Errors, imports.Errors)
		})
	}
}
