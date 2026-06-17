package dart

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLanguage_ParseFile(t *testing.T) {
	absTestFolder, _ := filepath.Abs(testFolder)

	tests := []struct {
		Name            string
		File            string
		ExpectedPackage string
		ExpectedRelPath string
	}{
		{
			Name:            "resolves package and rel path from pubspec.yaml",
			File:            filepath.Join("lib", "main.dart"),
			ExpectedPackage: "sample_app",
			ExpectedRelPath: filepath.Join("lib", "main.dart"),
		},
		{
			Name:            "resolves nested files relative to pubspec.yaml",
			File:            filepath.Join("lib", "src", "bar.dart"),
			ExpectedPackage: "sample_app",
			ExpectedRelPath: filepath.Join("lib", "src", "bar.dart"),
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
			a.Equal(tt.ExpectedPackage, file.Package)
			a.Equal(tt.ExpectedRelPath, file.RelPath)
		})
	}
}

func TestLanguage_ResolvePath(t *testing.T) {
	absTestFolder, _ := filepath.Abs(testFolder)
	libDir := filepath.Join(absTestFolder, "lib")

	tests := []struct {
		Name       string
		Unresolved string
		Dir        string
		Expected   string
	}{
		{
			Name:       "ignores dart sdk imports",
			Unresolved: "dart:async",
			Dir:        libDir,
			Expected:   "",
		},
		{
			Name:       "ignores third party package imports",
			Unresolved: "package:other_pkg/other.dart",
			Dir:        libDir,
			Expected:   "",
		},
		{
			Name:       "resolves own package imports against lib",
			Unresolved: "package:sample_app/foo.dart",
			Dir:        libDir,
			Expected:   filepath.Join(libDir, "foo.dart"),
		},
		{
			Name:       "resolves relative imports",
			Unresolved: "src/bar.dart",
			Dir:        libDir,
			Expected:   filepath.Join(libDir, "src", "bar.dart"),
		},
		{
			Name:       "returns empty for non existing files",
			Unresolved: "does/not/exist.dart",
			Dir:        libDir,
			Expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_lang, err := MakeDartLanguage(nil)
			a.NoError(err)
			lang := _lang.(*Language)

			resolved, err := lang.ResolvePath(tt.Unresolved, tt.Dir)
			a.NoError(err)
			a.Equal(tt.Expected, resolved)
		})
	}
}
