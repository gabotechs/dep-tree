package python

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/utils"
)

type InitModuleResult struct {
	Path        string
	PythonFiles []string
}

type DirectoryResult struct {
	Path        string
	PythonFiles []string
}

type FileResult struct {
	Path string
}

type ResolveResult struct {
	InitModule *InitModuleResult
	Directory  *DirectoryResult
	File       *FileResult
}

// resolveFromSlicesAndSearchPath returns multiple valid resolved paths.
func resolveFromSlicesAndSearchPath(searchPath string, slices []string) *ResolveResult {
	fullFileOrDir := path.Join(append([]string{searchPath}, slices...)...)

	if utils.FileExists(fullFileOrDir + ".py") {
		abs, _ := filepath.Abs(fullFileOrDir + ".py")
		return &ResolveResult{File: &FileResult{Path: abs}}
	}

	if result, err := os.ReadDir(fullFileOrDir); err == nil {
		var pythonFiles []string
		for _, entry := range result {
			if strings.HasSuffix(entry.Name(), ".py") {
				pythonFiles = append(pythonFiles, path.Join(fullFileOrDir, entry.Name()))
			}
		}
		initFile := path.Join(fullFileOrDir, "__init__.py")
		if utils.FileExists(initFile) {
			abs, _ := filepath.Abs(initFile)
			return &ResolveResult{InitModule: &InitModuleResult{
				Path:        abs,
				PythonFiles: pythonFiles,
			}}
		}
		abs, _ := filepath.Abs(fullFileOrDir)
		return &ResolveResult{Directory: &DirectoryResult{
			PythonFiles: pythonFiles,
			Path:        abs,
		}}
	}
	return nil
}

// ResolveRelative cannot return an empty []string, unless an error happened.
//
// In contrary to ResolveAbsolute, this method can return an error as a relative import is
// always expected to be found.
func ResolveRelative(slices []string, dir string, stepsBack int) (*ResolveResult, error) {
	var back []string
	for i := 0; i < stepsBack; i++ {
		back = append(back, "..")
	}
	searchPath := path.Join(append([]string{dir}, back...)...)
	result := resolveFromSlicesAndSearchPath(searchPath, slices)
	if result == nil {
		return nil, fmt.Errorf(
			"could not resolve relative import from %s to %s/%s",
			dir,
			searchPath,
			strings.Join(slices, "/"),
		)
	}
	return result, nil
}

// ResolveAbsolute never fails, if nothing is found it just returns nil.
//
// This is fine because we assume that an un-resolved absolute import is pointing to
// a library or something like that, so no need to take it into account.
func (l *Language) ResolveAbsolute(slices []string) *ResolveResult {
	searchPaths := l.PythonPath

	for _, searchPath := range searchPaths {
		result := resolveFromSlicesAndSearchPath(searchPath, slices)
		if result != nil {
			return result
		}
	}

	return nil
}
