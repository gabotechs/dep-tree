package golang

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	absPath, _ := filepath.Abs(".")

	tests := []struct {
		Name    string
		File    string
		RelPath string
		AbsPath string
	}{
		{
			Name:    "Can lookup a path given it's relative path",
			File:    "./language.go",
			RelPath: "internal/golang/language.go",
			AbsPath: filepath.Join(absPath, "language.go"),
		},
		{
			Name:    "Can lookup a path given it's absolute path",
			File:    filepath.Join(absPath, "language.go"),
			RelPath: "internal/golang/language.go",
			AbsPath: filepath.Join(absPath, "language.go"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			lang, err := NewLanguage(".", &Config{})
			a.NoError(err)
			file, err := lang.ParseFile(tt.File)
			a.NoError(err)
			a.Equal(tt.RelPath, file.RelPath)
			a.Equal(tt.AbsPath, file.AbsPath)
		})
	}
}
