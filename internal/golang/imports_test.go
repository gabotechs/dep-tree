package golang

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"

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
		{
			Name: "exports.go",
			Expected: [][2]string{
				{"Language", "internal/golang/language.go"},
				{"FileInfo", "internal/language/language.go"},
				{"File", "internal/golang/package.go"},
				{"ExportsResult", "internal/language/language.go"},
				{"ExportSymbol", "internal/language/language.go"},
				{"ExportEntry", "internal/language/language.go"},
			},
		},
		{
			Name: "package.go",
			Expected: [][2]string{
				{"Cached1In1OutErr", "internal/utils/cached.go"},
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

func Test_ImportStmt(t *testing.T) {
	tests := []struct {
		Name    string
		IsLocal bool
		RelPath string
		Alias   string
	}{
		{
			Name:    "github.com/stretchr/testify/require",
			IsLocal: false,
			RelPath: "",
			Alias:   "require",
		},
		{
			Name:    "github.com/gabotechs/dep-tree",
			IsLocal: false,
			RelPath: "",
			Alias:   "tree",
		},
		{
			Name:    "dep_tree github.com/gabotechs/dep-tree",
			IsLocal: false,
			RelPath: "",
			Alias:   "dep_tree",
		},
		{
			Name:    "github.com/gabotechs/dep-tree/internal/golang",
			IsLocal: true,
			RelPath: "internal/golang",
			Alias:   "golang",
		},
		{
			Name:    "go github.com/gabotechs/dep-tree/internal/golang",
			IsLocal: true,
			RelPath: "internal/golang",
			Alias:   "go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			lang, err := NewLanguage(".", &Config{})
			a.NoError(err)
			nameSlices := strings.Split(tt.Name, " ")
			var importStmt ImportStmt
			if len(nameSlices) == 1 {
				importStmt.importPath = nameSlices[0]
			} else {
				importStmt.importName = nameSlices[0]
				importStmt.importPath = nameSlices[1]
			}
			a.Equal(tt.IsLocal, importStmt.IsLocal(lang.GoMod.Module))
			a.Equal(tt.Alias, importStmt.Alias())
			a.Equal(tt.RelPath, importStmt.RelPath(lang.GoMod.Module))
		})
	}
}
