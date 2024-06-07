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
			ExpectedSymbol:      "PackagesInDir",
			ExpectedFileAbsPath: filepath.Join(absPath, "package.go"),
		},
		{
			Name:                "_packagesInDir function on this package",
			Path:                ".",
			ExpectedSymbol:      "_packagesInDir",
			ExpectedFileAbsPath: filepath.Join(absPath, "package.go"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			result, err := PackagesInDir(tt.Path)
			a.NoError(err)
			var pkg Package
			found := false
			for _, pkg = range result {
				if _, ok := pkg.SymbolToFile[tt.ExpectedSymbol]; ok {
					found = true
					break
				}
			}
			a.Equal(true, found)

			a.Equal(
				tt.ExpectedFileAbsPath,
				pkg.SymbolToFile[tt.ExpectedSymbol].AbsPath,
			)
		})
	}
}
