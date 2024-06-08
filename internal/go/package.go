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
	Package *Package
	AbsPath string
}

func _newFile(path string) (*File, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	absDir := filepath.Dir(absPath)
	pkgs, err := PackagesInDir(absDir)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		if file, ok := pkg.AbsPathToFile[absPath]; ok {
			return file, nil
		}
	}

	return nil, fmt.Errorf("could not find file %s in any of the loaded packages", absPath)
}

var NewFile = utils.Cached1In1OutErr(_newFile)

type Package struct {
	*ast.Package
	SymbolToFile  map[string]*File
	AbsPathToFile map[string]*File
}

func _packagesInDir(dir string) ([]*Package, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, absDir, nil, 0)
	if err != nil {
		return nil, err
	}

	result := make([]*Package, len(pkgs))
	i := 0
	for _, pkg := range pkgs {
		filePkg := Package{
			Package:       pkg,
			SymbolToFile:  make(map[string]*File),
			AbsPathToFile: make(map[string]*File),
		}
		for absFilePath, file := range pkg.Files {
			f := File{
				File:    file,
				AbsPath: absFilePath,
				Package: &filePkg,
			}
			filePkg.AbsPathToFile[absFilePath] = &f
			for name := range file.Scope.Objects {
				filePkg.SymbolToFile[name] = &f
			}
		}
		result[i] = &filePkg
		i += 1
	}
	return result, nil
}

var PackagesInDir = utils.Cached1In1OutErr(_packagesInDir)
