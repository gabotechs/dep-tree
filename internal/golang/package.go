package golang

import (
	"fmt"
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

func _NewPackageFromDir(dir string) (*Package, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, absDir, nil, 0)
	if err != nil {
		return nil, err
	}

	var pkg *ast.Package
	for _, pkg = range pkgs {
		// There should only be one entry in `pkgs`, I don't know how there would
		// be multiple.
	}
	if pkg == nil {
		return nil, fmt.Errorf("could not find any packages in directory %s", dir)
	}

	exports := make(map[string]*File)
	for absFileName, file := range pkg.Files {
		for name := range file.Scope.Objects {
			exports[name] = &File{
				File:    file,
				AbsPath: absFileName,
			}
		}
	}
	return &Package{
		Package:      pkg,
		SymbolToFile: exports,
	}, nil
}

var NewPackageFromDir = utils.Cached1In1OutErr(_NewPackageFromDir)
