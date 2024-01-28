package python

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/utils"
)

type InitModuleResult struct {
	Path        string
	PythonFiles []string
}

func (i *InitModuleResult) fileMap() map[string]string {
	availableFiles := map[string]string{}
	for _, pythonFile := range i.PythonFiles {
		for _, ext := range Extensions {
			fileName := filepath.Base(pythonFile)
			if strings.HasSuffix(fileName, ext) {
				availableFiles[strings.TrimSuffix(fileName, "."+ext)] = pythonFile
			}
		}
	}
	return availableFiles
}

type DirectoryResult struct {
	Path        string
	PythonFiles []string
}

func (d *DirectoryResult) fileMap() map[string]string {
	availableFiles := map[string]string{}
	for _, pythonFile := range d.PythonFiles {
		availableFiles[strings.TrimSuffix(filepath.Base(pythonFile), ".py")] = pythonFile
	}
	return availableFiles
}

type FileResult struct {
	Path string
}

type ResolveResult struct {
	InitModule *InitModuleResult
	Directory  *DirectoryResult
	File       *FileResult
}

func _pythonFilesInDir(dir string) []string {
	result, _ := os.ReadDir(dir)
	var pythonFiles []string
	for _, entry := range result {
		for _, ext := range Extensions {
			if strings.HasSuffix(entry.Name(), "."+ext) {
				pythonFiles = append(pythonFiles, filepath.Join(dir, entry.Name()))
			}
		}
	}
	return pythonFiles
}

var pythonFilesInDir = utils.Cached1In1Out(_pythonFilesInDir)

// resolveFromSlicesAndSearchPath returns multiple valid resolved paths.
func resolveFromSlicesAndSearchPath(searchPath string, slices []string) *ResolveResult {
	fullFileOrDir := filepath.Join(append([]string{searchPath}, slices...)...)

	// If there is a Python file, we are done.
	for _, ext := range Extensions {
		if utils.FileExists(fullFileOrDir + "." + ext) {
			abs, _ := filepath.Abs(fullFileOrDir + "." + ext)
			return &ResolveResult{File: &FileResult{Path: abs}}
		}
	}

	// If there was not a Python file, it should be a dir.
	if !utils.DirExists(fullFileOrDir) {
		return nil
	}

	pythonFiles := pythonFilesInDir(fullFileOrDir)
	// If there is an __init__.py file, we must be referring to that one.
	initFile := filepath.Join(fullFileOrDir, "__init__.py")
	if utils.FileExists(initFile) {
		abs, _ := filepath.Abs(initFile)
		return &ResolveResult{InitModule: &InitModuleResult{
			Path:        abs,
			PythonFiles: pythonFiles,
		}}
	}
	// Otherwise, the whole folder is being imported.
	abs, _ := filepath.Abs(fullFileOrDir)
	return &ResolveResult{Directory: &DirectoryResult{
		PythonFiles: pythonFiles,
		Path:        abs,
	}}
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
	searchPath := filepath.Join(append([]string{dir}, back...)...)
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

var rootFiles = []string{
	"pyproject.toml",
	"setup.py",
	"poetry.toml",
	"poetry.lock",
	"requirements.txt",
	".pylintrc",
	".git/index",
}

func _findClosestDirWithRootFile(searchPath string) string {
	for _, rootFile := range rootFiles {
		if utils.FileExists(filepath.Join(searchPath, rootFile)) {
			return searchPath
		}
	}
	nextSearchPath := filepath.Dir(searchPath)

	if nextSearchPath != searchPath {
		return _findClosestDirWithRootFile(nextSearchPath)
	} else {
		return ""
	}
}

var findClosestDirWithRootFile = utils.Cached1In1Out(_findClosestDirWithRootFile)

// ResolveAbsolute never fails, if nothing is found it just returns nil.
//
// This is fine because we assume that an un-resolved absolute import is pointing to
// a library or something like that, so no need to take it into account.
func (l *Language) ResolveAbsolute(slices []string, currDir string) *ResolveResult {
	searchPaths := []string{currDir}
	dirWithRootFile := findClosestDirWithRootFile(currDir)
	if dirWithRootFile != "" {
		searchPaths = append(searchPaths, dirWithRootFile)
	}
	if l.cfg != nil {
		searchPaths = append(searchPaths, l.cfg.PythonPath...)
	}

	for _, searchPath := range searchPaths {
		result := resolveFromSlicesAndSearchPath(searchPath, slices)
		if result != nil {
			return result
		}
	}

	return nil
}
