package js

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"dep-tree/internal/utils"
)

type ResolveCacheKey string

func makeResolveCacheKey(unresolved string, dir string) ResolveCacheKey {
	return ResolveCacheKey(unresolved + dir)
}

func (l *Language) ResolvePath(ctx context.Context, unresolved string, dir string) (context.Context, string, error) {
	cacheKey := makeResolveCacheKey(unresolved, dir)
	if cached, ok := ctx.Value(cacheKey).(string); ok {
		return ctx, cached, nil
	} else {
		resolved, err := l._uncachedResolvePath(unresolved, dir)
		if err != nil {
			return ctx, "", err
		}
		ctx = context.WithValue(ctx, cacheKey, resolved)
		return ctx, resolved, nil
	}
}

// ResolvePath resolves an unresolved import based on the dir where the import was executed.
func (l *Language) _uncachedResolvePath(unresolved string, dir string) (string, error) {
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

	// 2. If is imported from baseUrl.
	baseUrl := l.TsConfig.CompilerOptions.BaseUrl
	importFromBaseUrl := path.Join(l.ProjectRoot, baseUrl, unresolved)
	absPath = getFileAbsPath(importFromBaseUrl)
	if absPath != "" {
		return absPath, nil
	}

	// 3. If imported from a path override.
	pathOverrides := l.TsConfig.CompilerOptions.Paths
	if pathOverrides == nil {
		return absPath, nil
	}
	for pathOverride, searchPaths := range pathOverrides {
		pathOverride = strings.ReplaceAll(pathOverride, "*", "")
		if strings.HasPrefix(unresolved, pathOverride) {
			for _, searchPath := range searchPaths {
				searchPath = strings.ReplaceAll(searchPath, "*", "")
				newImportFrom := strings.ReplaceAll(unresolved, pathOverride, searchPath)
				importFromBaseUrlAndPaths := path.Join(l.ProjectRoot, baseUrl, newImportFrom)
				absPath = getFileAbsPath(importFromBaseUrlAndPaths)
				if absPath != "" {
					return absPath, nil
				}
			}
			return absPath, fmt.Errorf("import '%s' was matched to path '%s' in tscofing's paths option, but the resolved path did not match an existing file", unresolved, pathOverride)
		}
	}
	return absPath, nil
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
