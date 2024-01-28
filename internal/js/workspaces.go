package js

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/utils"
)

type Workspaces struct {
	// ws is a map from packageJson name to absolute path.
	ws map[string]*packageJson
}

func allDirsWithAPackageJson(start string) ([]string, error) {
	dir, err := os.ReadDir(start)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, entry := range dir {
		if entry.IsDir() {
			more, err := allDirsWithAPackageJson(filepath.Join(start, entry.Name()))
			if err != nil {
				return nil, err
			}
			result = append(result, more...)
		} else if entry.Name() == packageJsonFile {
			result = append(result, start)
		}
	}
	return result, nil
}

func searchFirstPackageJsonWithWorkspaces(searchPath string) (*packageJson, error) {
	packageJsonPath := filepath.Join(searchPath, packageJsonFile)
	if utils.FileExists(packageJsonPath) {
		result, err := readPackageJson(packageJsonPath)
		if err != nil {
			return nil, err
		}
		if len(result.workspaces()) > 0 {
			return result, nil
		}
	}
	nextSearchPath := filepath.Dir(searchPath)
	if nextSearchPath != searchPath {
		return searchFirstPackageJsonWithWorkspaces(nextSearchPath)
	}
	return nil, nil
}

var NewWorkspaces = utils.Cached1In1OutErr(func(searchPath string) (*Workspaces, error) {
	searchPath, err := filepath.Abs(searchPath)
	if err != nil {
		return nil, err
	}
	rootPackageJson, err := searchFirstPackageJsonWithWorkspaces(searchPath)
	if err != nil {
		return nil, err
	}
	if rootPackageJson == nil {
		return nil, nil
	}
	dirsWithAPackageJson, err := allDirsWithAPackageJson(rootPackageJson.absPath)
	workspacesMap := map[string]*packageJson{}

	for _, dir := range dirsWithAPackageJson {
		for _, ws := range rootPackageJson.workspaces() {
			rel, _ := filepath.Rel(rootPackageJson.absPath, dir)
			match, err := utils.GlobstarMatch(ws, rel)
			if err != nil {
				return nil, err
			}
			if match {
				pkgJson, err := readPackageJson(dir)
				if err != nil {
					return nil, err
				}
				workspacesMap[pkgJson.Name] = pkgJson
			}
		}
	}
	return &Workspaces{ws: workspacesMap}, nil
})

func (w *Workspaces) ResolveFromWorkspaces(unresolved string) (string, error) {
	if w == nil {
		return "", nil
	}
	slices := strings.SplitN(unresolved, "/", 2)
	firstSlice := slices[0]
	rest := ""
	if len(slices) > 1 {
		rest = slices[1]
	}
	var pkgJson *packageJson
	for {
		entry, ok := w.ws[firstSlice]
		//nolint:gocritic
		if ok {
			pkgJson = entry
			break
		} else if rest == "" {
			return "", nil
		} else {
			slices = strings.SplitN(rest, "/", 2)
			firstSlice += "/" + slices[0]
			rest = ""
			if len(slices) > 1 {
				rest = slices[1]
			}
		}
	}
	var fullPath string
	if rest == "" {
		fullPath = pkgJson.index()
		if fullPath == "" {
			return "", fmt.Errorf("workspace '%s' has no index file", pkgJson.absPath)
		}
	} else {
		fullPath = filepath.Join(pkgJson.absPath, rest)
	}
	result := getFileAbsPath(fullPath)
	if result == "" {
		return "", fmt.Errorf("path '%s' matched workspace '%s', but no file '%s' does not exist", unresolved, firstSlice, fullPath)
	}
	return result, nil
}
