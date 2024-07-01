package js

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/tailscale/hujson"
)

const tsConfigFile = "tsconfig.json"

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
	tsConfig.path = filepath.Dir(filePath)
	return tsConfig, err
}

func (t *TsConfig) ResolveFromBaseUrl(unresolved string) string {
	return filepath.Join(t.path, t.CompilerOptions.BaseUrl, unresolved)
}

func (t *TsConfig) ResolveFromPaths(unresolved string) []string {
	var candidates []string

	for pathOverride, searchPaths := range t.CompilerOptions.Paths {
		pathOverride = strings.ReplaceAll(pathOverride, "*", "")
		if strings.HasPrefix(unresolved, pathOverride) {
			for _, searchPath := range searchPaths {
				searchPath = strings.ReplaceAll(searchPath, "*", "")
				newImportFrom := strings.ReplaceAll(unresolved, pathOverride, searchPath)
				importFromBaseUrlAndPaths := filepath.Join(t.path, t.CompilerOptions.BaseUrl, newImportFrom)
				candidates = append(candidates, importFromBaseUrlAndPaths)
			}
		}
	}
	return candidates
}
