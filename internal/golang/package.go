package golang

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/utils"
)

type File struct {
	*ast.File
	AbsPath string
}

type Package struct {
	*ast.Package
	SymbolToFile map[string]*File
}

func _packagesInDir(dir string) ([]Package, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, absDir, nil, 0)
	if err != nil {
		return nil, err
	}

	result := make([]Package, len(pkgs))
	i := 0
	for _, pkg := range pkgs {
		exports := make(map[string]*File)
		for absFileName, file := range pkg.Files {
			for name := range file.Scope.Objects {
				exports[name] = &File{
					File:    file,
					AbsPath: absFileName,
				}
			}
		}
		result[i] = Package{
			Package:      pkg,
			SymbolToFile: exports,
		}
		i += 1
	}
	return result, nil
}

var PackagesInDir = utils.Cached1In1OutErr(_packagesInDir)
