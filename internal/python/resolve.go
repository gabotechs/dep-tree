package python

import (
	"fmt"
	"path"
	"strings"

	"dep-tree/internal/utils"
)

type ResolveResult struct {
	InitModule string
	Directory  string
	File       string
}

// resolveFromSlicesAndSearchPath returns multiple valid resolved paths.
func resolveFromSlicesAndSearchPath(searchPath string, slices []string) *ResolveResult {
	fullFileOrDir := path.Join(append([]string{searchPath}, slices...)...)

	if utils.FileExists(fullFileOrDir + ".py") {
		return &ResolveResult{File: fullFileOrDir + ".py"}
	}

	if utils.DirExists(fullFileOrDir) {
		initFile := path.Join(fullFileOrDir, "__init__.py")
		if utils.FileExists(initFile) {
			return &ResolveResult{InitModule: initFile}
		} else {
			return &ResolveResult{Directory: fullFileOrDir}
		}
	}
	return nil
}

// ResolveRelative cannot return an empty []string, unless an error happened.
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

func (l *Language) ResolveAbsolute(slices []string) (*ResolveResult, error) {
	searchPaths := l.PythonPath

	for _, searchPath := range searchPaths {
		result := resolveFromSlicesAndSearchPath(searchPath, slices)
		if result != nil {
			return result, nil
		}
	}

	return nil, nil
}
