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

	file, err := NewFile(path)
	if err != nil {
		return nil, err
	}

	relPath, _ := filepath.Rel(l.Root.AbsDir, absPath)

	return &language.FileInfo{
		Content: file,
		AbsPath: absPath,
		RelPath: relPath,
		Package: file.Package.Name,
		Size:    file.TokenFile.Size(),
		Loc:     file.TokenFile.LineCount(),
	}, nil
}

var findClosestDirWithRootFile = utils.MakeCachedFindClosestDirWithRootFile([]string{
	// NOTE: for now, only support projects that contain a go.mod file.
	"go.mod",
})
