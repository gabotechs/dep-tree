package js

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/tailscale/hujson"
)

type CompilerOptions struct {
	BaseUrl string              `json:"baseUrl,omitempty"`
	Paths   map[string][]string `json:"paths,omitempty"`
}

type TsConfig struct {
	path            string
	CompilerOptions CompilerOptions `json:"compilerOptions,omitempty"`
}

func ParseTsConfig(filePath string) (TsConfig, error) {
	var tsConfig TsConfig
	data, err := os.ReadFile(filePath)
	if err != nil {
		return TsConfig{}, err
	}
	standard, err := hujson.Standardize(data)
	if err != nil {
		return TsConfig{}, err
	}
	err = json.Unmarshal(standard, &tsConfig)
	if err != nil {
		return TsConfig{}, err
	}
	tsConfig.path = path.Dir(filePath)
	return tsConfig, err
}

func (t *TsConfig) ResolveFromBaseUrl(unresolved string) string {
	baseUrl := t.CompilerOptions.BaseUrl
	importFromBaseUrl := path.Join(t.path, baseUrl, unresolved)
	return getFileAbsPath(importFromBaseUrl)
}

func (t *TsConfig) ResolveFromPaths(unresolved string) (string, error) {
	var failedMatches []string
	for pathOverride, searchPaths := range t.CompilerOptions.Paths {
		pathOverride = strings.ReplaceAll(pathOverride, "*", "")
		if strings.HasPrefix(unresolved, pathOverride) {
			for _, searchPath := range searchPaths {
				searchPath = strings.ReplaceAll(searchPath, "*", "")
				newImportFrom := strings.ReplaceAll(unresolved, pathOverride, searchPath)
				importFromBaseUrlAndPaths := path.Join(t.path, t.CompilerOptions.BaseUrl, newImportFrom)
				absPath := getFileAbsPath(importFromBaseUrlAndPaths)
				if absPath != "" {
					return absPath, nil
				}
			}
			failedMatches = append(failedMatches, pathOverride)
		}
	}
	if failedMatches != nil {
		return "", fmt.Errorf("import '%s' was matched to path '%s' in tscofing's paths option, but the resolved path did not match an existing file", unresolved, strings.Join(failedMatches, "', '"))
	} else {
		return "", nil
	}
}
