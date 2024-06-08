package golang_test

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"

	golang "github.com/gabotechs/dep-tree/internal/go"
	. "github.com/gabotechs/dep-tree/internal/language"
	"github.com/stretchr/testify/require"
)

type dummy = Language

func TestImports(t *testing.T) {
	tests := []struct {
		Name     string
		Expected [][2]string
	}{
		{
			Name: "imports.go",
			Expected: [][2]string{
				{"SymbolsImport", "internal/language/language.go"},
				{"Package", "internal/go/package.go"},
				{"PackagesInDir", "internal/go/package.go"},
				{"Language", "internal/go/language.go"},
				{"ImportsResult", "internal/language/language.go"},
				{"FileInfo", "internal/language/language.go"},
				{"File", "internal/go/package.go"},
			},
		},
		{
			Name: "exports.go",
			Expected: [][2]string{
				{"Language", "internal/go/language.go"},
				{"FileInfo", "internal/language/language.go"},
				{"File", "internal/go/package.go"},
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
		{
			Name: "imports_test.go",
			Expected: [][2]string{
				{"NewLanguage", "internal/go/language.go"},
				{"ImportStmt", "internal/go/imports.go"},
				{"Config", "internal/go/config.go"},
				{"Language", "internal/language/language.go"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			lang, err := golang.NewLanguage(".", &golang.Config{})
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
				return expected[i][0] > expected[j][0]
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
			lang, err := golang.NewLanguage(".", &golang.Config{})
			a.NoError(err)
			nameSlices := strings.Split(tt.Name, " ")
			var importStmt golang.ImportStmt
			if len(nameSlices) == 1 {
				importStmt.ImportPath = nameSlices[0]
			} else {
				importStmt.ImportName = nameSlices[0]
				importStmt.ImportPath = nameSlices[1]
			}
			a.Equal(tt.IsLocal, importStmt.IsLocal(lang.GoMod.Module))
			a.Equal(tt.Alias, importStmt.Alias())
			a.Equal(tt.RelPath, importStmt.RelPath(lang.GoMod.Module))
		})
	}
}
