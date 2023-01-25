package js

import (
	"encoding/json"
	"os"

	"github.com/tailscale/hujson"
)

type CompilerOptions struct {
	BaseUrl string              `json:"baseUrl,omitempty"`
	Paths   map[string][]string `json:"paths,omitempty"`
}

type TsConfig struct {
	CompilerOptions CompilerOptions `json:"compilerOptions,omitempty"`
}

func ParseTsConfig(path string) (TsConfig, error) {
	var tsConfig TsConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return TsConfig{}, err
	}
	standard, err := hujson.Standardize(data)
	if err != nil {
		return TsConfig{}, err
	}
	err = json.Unmarshal(standard, &tsConfig)
	return tsConfig, err
}
