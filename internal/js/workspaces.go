package js

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/utils"
)

type WorkspaceEntry struct {
	absPath string
	main    string
}

func (w *WorkspaceEntry) index() string {
	// Independently of what the package.json `main` says, let's
	// always try first the `src/index.[js|ts|jsx|tsx]` file.
	fullPath := getFileAbsPath(path.Join(w.absPath, "src"))
	if fullPath != "" {
		return fullPath
	}
	// Then, if a `main` property is declared in the package.json, follow it.
	if w.main != "" {
		fullPath = getFileAbsPath(path.Join(w.absPath, w.main))
		if fullPath != "" {
			return fullPath
		}
	}
	// Then, as a last resource, check if there is an `index.[js|ts|jsx|tsx]`
	// file in the root of the project.
	return getFileAbsPath(w.absPath)
}

type Workspaces struct {
	// ws is a map from packageJson name to absolute path.
	ws map[string]WorkspaceEntry
}

type partialPackageJson struct {
	path       string
	Main       string   `json:"main,omitempty"`
	Name       string   `json:"name"`
	Workspaces []string `json:"workspaces"`
}

func searchFirstPackageJsonWithWorkspaces(searchPath string) (*partialPackageJson, error) {
	if len(searchPath) < 2 {
		return nil, nil
	}
	packageJsonPath := path.Join(searchPath, "package.json")
	if utils.FileExists(packageJsonPath) {
		var result partialPackageJson
		content, err := os.ReadFile(packageJsonPath)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(content, &result)
		if err != nil {
			return nil, err
		}
		if len(result.Workspaces) > 0 {
			result.path = searchPath
			return &result, nil
		} else {
			return searchFirstPackageJsonWithWorkspaces(path.Dir(searchPath))
		}
	} else {
		return searchFirstPackageJsonWithWorkspaces(path.Dir(searchPath))
	}
}

func allDirsWithAPackageJson(start string) ([]string, error) {
	dir, err := os.ReadDir(start)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, entry := range dir {
		if entry.IsDir() {
			more, err := allDirsWithAPackageJson(path.Join(start, entry.Name()))
			if err != nil {
				return nil, err
			}
			result = append(result, more...)
		} else if entry.Name() == "package.json" {
			result = append(result, start)
		}
	}
	return result, nil
}

func NewWorkspaces(searchPath string) (*Workspaces, error) {
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
	dirsWithAPackageJson, err := allDirsWithAPackageJson(rootPackageJson.path)
	workspacesMap := map[string]WorkspaceEntry{}

	for _, dir := range dirsWithAPackageJson {
		for _, ws := range rootPackageJson.Workspaces {
			rel, _ := filepath.Rel(rootPackageJson.path, dir)
			match, err := utils.GlobstarMatch(ws, rel)
			if err != nil {
				return nil, err
			}
			if match {
				packageJsonPath := path.Join(dir, "package.json")
				content, err := os.ReadFile(packageJsonPath)
				if err != nil {
					return nil, err
				}
				var packageJson partialPackageJson
				err = json.Unmarshal(content, &packageJson)
				if err != nil {
					return nil, fmt.Errorf("cannot parse %s: %w", packageJsonPath, err)
				}
				workspacesMap[packageJson.Name] = WorkspaceEntry{
					absPath: dir,
					main:    packageJson.Main,
				}
			}
		}
	}
	return &Workspaces{ws: workspacesMap}, nil
}

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
	var ws WorkspaceEntry
	for {
		entry, ok := w.ws[firstSlice]
		//nolint:gocritic
		if ok {
			ws = entry
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
		fullPath = ws.index()
		if fullPath == "" {
			return "", fmt.Errorf("workspace '%s' has no index file", ws.absPath)
		}
	} else {
		fullPath = path.Join(ws.absPath, rest)
	}
	result := getFileAbsPath(fullPath)
	if result == "" {
		return "", fmt.Errorf("path '%s' matched workspace '%s', but no file '%s' does not exist", unresolved, firstSlice, fullPath)
	}
	return result, nil
}
