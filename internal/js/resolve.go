package js

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/utils"
)

// ResolvePath resolves an unresolved import based on the dir where the import was executed.
//
//nolint:gocyclo
func (l *Language) ResolvePath(unresolved string, dir string) (string, error) {
	absPath := ""

	if len(unresolved) == 0 {
		return "", errors.New("import path cannot be empty")
	} else if len(unresolved) == 1 {
		return "", fmt.Errorf("invalid import path %s", unresolved)
	}

	// 1. If import is relative.
	if unresolved[0] == '.' && (unresolved[1] == '/' || unresolved[1] == '.') {
		absPath = getFileAbsPath(path.Join(dir, unresolved))
		if absPath == "" {
			return absPath, fmt.Errorf("could not perform relative import for '%s' because the file or dir was not found", unresolved)
		}
		return absPath, nil
	}

	tsConfig, _, err := findPackageJson(dir)
	if err != nil {
		return "", err
	}

	// 2. If is imported from a workspace.
	if l.Cfg == nil || l.Cfg.Workspaces {
		absPath, err = l.Workspaces.ResolveFromWorkspaces(unresolved)
		if absPath != "" || err != nil {
			return absPath, err
		}
	}

	// 3. If is imported from baseUrl.
	absPath = tsConfig.ResolveFromBaseUrl(unresolved)
	if absPath != "" {
		return absPath, nil
	}

	// 4. If imported from a path override.
	if l.Cfg == nil || l.Cfg.TsConfigPaths {
		absPath, err = tsConfig.ResolveFromPaths(unresolved)
		if err != nil {
			return "", err
		}
		if absPath != "" {
			return absPath, nil
		}
	}
	return "", nil
}

func retrieveWithExt(absPath string) string {
	for _, ext := range Extensions {
		if strings.HasSuffix(absPath, "."+ext) {
			absPath = absPath[0 : len(absPath)-len("."+ext)]
		}
	}
	for _, ext := range Extensions {
		withExtPath := absPath + "." + ext
		if utils.FileExists(withExtPath) {
			return withExtPath
		}
	}
	return ""
}

func getFileAbsPath(id string) string {
	absPath, err := filepath.Abs(id)
	switch {
	case err != nil:
		return ""
	case utils.DirExists(id):
		return retrieveWithExt(path.Join(absPath, "index"))
	default:
		return retrieveWithExt(absPath)
	}
}
