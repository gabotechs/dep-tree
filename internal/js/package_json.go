package js

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/utils"
)

const packageJsonFile = "package.json"

type packageJson struct {
	absPath    string
	Main       string      `json:"main,omitempty"`
	Name       string      `json:"name"`
	Workspaces interface{} `json:"workspaces"`
}

var readPackageJson = utils.Cached1In2Out(func(path string) (*packageJson, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	fullPath := ""
	dir := ""
	if filepath.Base(path) != packageJsonFile {
		fullPath = filepath.Join(path, packageJsonFile)
		dir = path
	} else {
		fullPath = path
		dir = filepath.Dir(path)
	}
	var result packageJson
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing %q: %w", fullPath, err)
	}
	result.absPath = dir
	return &result, nil
})

func castAnyArray[T any](arr []any) []T {
	result := make([]T, len(arr))
	for i, el := range arr {
		result[i] = el.(T)
	}
	return result
}

func (p *packageJson) workspaces() []string {
	switch v := p.Workspaces.(type) {
	case []any:
		return castAnyArray[string](v)
	case map[string]any:
		if packages, ok := v["packages"]; ok {
			if vv, ok := packages.([]any); ok {
				return castAnyArray[string](vv)
			}
		}
	}
	return []string{}
}

func (p *packageJson) index() string {
	// Independently of what the package.json `main` says, let's
	// always try first the `src/index.[js|ts|jsx|tsx]` file.
	fullPath := getFileAbsPath(filepath.Join(p.absPath, "src"))
	if fullPath != "" {
		return fullPath
	}
	// Then, if a `main` property is declared in the package.json, follow it.
	if p.Main != "" {
		fullPath = getFileAbsPath(filepath.Join(p.absPath, p.Main))
		if fullPath != "" {
			return fullPath
		}
	}
	// Then, as a last resource, check if there is an `index.[js|ts|jsx|tsx]`
	// file in the root of the project.
	return getFileAbsPath(p.absPath)
}
