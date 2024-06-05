package golang

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackage(t *testing.T) {
	absPath, _ := filepath.Abs(".")

	tests := []struct {
		Name                string
		Path                string
		ExpectedSymbol      string
		ExpectedFileAbsPath string
	}{
		{
			Name:                "File type on this package",
			Path:                ".",
			ExpectedSymbol:      "File",
			ExpectedFileAbsPath: filepath.Join(absPath, "package.go"),
		},
		{
			Name:                "NewPackageFromDir function on this package",
			Path:                ".",
			ExpectedSymbol:      "NewPackageFromDir",
			ExpectedFileAbsPath: filepath.Join(absPath, "package.go"),
		},
		{
			Name:                "_NewPackageFromDir function on this package",
			Path:                ".",
			ExpectedSymbol:      "_NewPackageFromDir",
			ExpectedFileAbsPath: filepath.Join(absPath, "package.go"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			result, err := NewPackageFromDir(tt.Path)
			a.NoError(err)

			a.Equal(
				tt.ExpectedFileAbsPath,
				result.SymbolToFile[tt.ExpectedSymbol].AbsPath,
			)
		})
	}
}
