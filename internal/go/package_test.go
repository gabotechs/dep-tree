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
		ExpectedPackage     string
	}{
		{
			Name:                "File type on this package",
			Path:                ".",
			ExpectedSymbol:      "File",
			ExpectedFileAbsPath: filepath.Join(absPath, "package.go"),
			ExpectedPackage:     "golang",
		},
		{
			Name:                "NewPackageFromDir function on this package",
			Path:                ".",
			ExpectedSymbol:      "PackagesInDir",
			ExpectedFileAbsPath: filepath.Join(absPath, "package.go"),
			ExpectedPackage:     "golang",
		},
		{
			Name:                "_packagesInDir function on this package",
			Path:                ".",
			ExpectedSymbol:      "_packagesInDir",
			ExpectedFileAbsPath: filepath.Join(absPath, "package.go"),
			ExpectedPackage:     "golang",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			result, err := PackagesInDir(tt.Path)
			a.NoError(err)
			var pkg *Package
			var file *File
			found := false
			for _, pkg = range result {
				var ok bool
				if file, ok = pkg.SymbolToFile[tt.ExpectedSymbol]; ok {
					found = true
					break
				}
			}
			a.Equal(true, found)

			a.Equal(
				tt.ExpectedFileAbsPath,
				pkg.SymbolToFile[tt.ExpectedSymbol].AbsPath,
			)

			a.Equal(tt.ExpectedPackage, file.Package.Name)
		})
	}
}

func TestFile(t *testing.T) {
	tests := []struct {
		Name        string
		PackageName string
	}{
		{
			Name:        "package.go",
			PackageName: "golang",
		},
		{
			Name:        "imports_test.go",
			PackageName: "golang_test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			file, err := NewFile(tt.Name)
			a.NoError(err)

			a.Equal(tt.PackageName, file.Package.Name)
			a.NotNil(file.Package)
			a.NotNil(file.TokenFile)
		})
	}
}
