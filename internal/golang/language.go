package golang

import (
	"fmt"
	"go/ast"
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/utils"
)

var Extensions = []string{
	"go",
}

type Language struct {
	Cfg      *Config
	GoMod    GoMod
	Root     utils.SourcesRoot
	Packages map[string]*ast.Package
}

func NewLanguage(dir string, cfg *Config) (*Language, error) {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	sourcesRoot := findClosestDirWithRootFile(absPath)
	if sourcesRoot == nil {
		return nil, fmt.Errorf("no go.mod file found in any parent directory of %s", dir)
	}
	goMod, err := ParseGoMod(filepath.Join(sourcesRoot.AbsDir, sourcesRoot.FoundFile))
	if err != nil {
		return nil, err
	}

	return &Language{
		Cfg:   cfg,
		GoMod: *goMod,
		Root:  *sourcesRoot,
	}, nil
}

func (l *Language) ParseFile(path string) (*language.FileInfo, error) {
	if l.Packages == nil {
		l.Packages = make(map[string]*ast.Package)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	absDir := filepath.Dir(absPath)

	pkgs, err := PackagesInDir(absDir)
	if err != nil {
		return nil, err
	}
	var file *ast.File
	var pkg Package
	for _, pkg = range pkgs {
		var ok bool
		if file, ok = pkg.Files[absPath]; ok {
			break
		}
	}
	if file == nil {
		return nil, fmt.Errorf("could not find file %s in any package in dir %s", absPath, absDir)
	}

	relPath, _ := filepath.Rel(l.Root.AbsDir, absPath)

	return &language.FileInfo{
		Content: &File{
			File:    file,
			AbsPath: absPath,
		},
		AbsPath: absPath,
		RelPath: relPath,
		Package: pkg.Name,
		Size:    int(file.FileEnd),
		Loc:     0, // TODO: I still don't know how to extract the LOC from an `ast.File` object.
	}, nil
}

var findClosestDirWithRootFile = utils.MakeCachedFindClosestDirWithRootFile([]string{
	// NOTE: for now, only support projects that contain a go.mod file.
	"go.mod",
})
