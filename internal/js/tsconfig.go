package js

import (
	"encoding/json"
	"os"
)

type TsConfig struct {
	CompilerOptions struct {
		BaseUrl string              `json:"baseUrl,omitempty"`
		Paths   map[string][]string `json:"paths,omitempty"`
	} `json:"compilerOptions,omitempty"`
}

func ParseTsConfig(path string) (TsConfig, error) {
	var tsConfig TsConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return TsConfig{}, err
	}
	err = json.Unmarshal(data, &tsConfig)
	if err != nil {
		return tsConfig, err
	}
	return tsConfig, nil
}
