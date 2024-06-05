package golang

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/gabotechs/dep-tree/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestImports(t *testing.T) {
	tests := []struct {
		Name     string
		Expected [][2]string
	}{
		{
			Name: "imports.go",
			Expected: [][2]string{
				{"SymbolsImport", "internal/language/language.go"},
				{"Package", "internal/golang/package.go"},
				{"NewPackageFromDir", "internal/golang/package.go"},
				{"Language", "internal/golang/language.go"},
				{"ImportsResult", "internal/language/language.go"},
				{"FileInfo", "internal/language/language.go"},
				{"File", "internal/golang/package.go"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			lang, err := NewLanguage(".", &Config{})
			a.NoError(err)
			file, err := lang.ParseFile(tt.Name)
			a.NoError(err)
			imports, err := lang.ParseImports(file)
			a.NoError(err)

			var actual [][2]string
			for _, imp := range imports.Imports {
				a.Equal(1, len(imp.Symbols))
				actual = append(actual, [2]string{imp.Symbols[0], imp.AbsPath})
			}

			var expected [][2]string
			for _, imp := range tt.Expected {
				expected = append(expected,
					[2]string{imp[0], filepath.Join(lang.Root.AbsDir, imp[1])},
				)
			}

			sort.Slice(actual, func(i, j int) bool {
				return actual[i][0] > actual[j][0]
			})
			sort.Slice(expected, func(i, j int) bool {
				return actual[i][0] > actual[j][0]
			})

			a.Equal(expected, actual)
		})
	}
}

func Test_importToPath(t *testing.T) {
	tests := []struct {
		Name     string
		Expected string
	}{
		{
			Name:     "github.com/stretchr/testify/require",
			Expected: "",
		},
		{
			Name:     "github.com/gabotechs/dep-tree",
			Expected: "",
		},
		{
			Name:     "github.com/gabotechs/dep-tree/internal/golang",
			Expected: "internal/golang",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			lang, err := NewLanguage(".", &Config{})
			a.NoError(err)
			result := lang.importToPath(tt.Name)
			a.Equal(tt.Expected, result)
			if result != "" {
				absDir := filepath.Join(lang.Root.AbsDir, result)
				a.Equal(true, utils.DirExists(absDir))
			}
		})
	}
}

func Test_importToAlias(t *testing.T) {
	tests := []struct {
		Name     string
		Expected string
	}{
		{
			Name:     "github.com/stretchr/testify/require",
			Expected: "require",
		},
		{
			Name:     "github.com/gabotechs/dep-tree",
			Expected: "tree",
		},
		{
			Name:     "github.com/gabotechs/dep-tree/internal/golang",
			Expected: "golang",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			a.Equal(tt.Expected, importToAlias(tt.Name))
		})
	}
}
