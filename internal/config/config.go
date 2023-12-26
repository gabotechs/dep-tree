package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/gabotechs/dep-tree/internal/check"
	"github.com/gabotechs/dep-tree/internal/js"
	"github.com/gabotechs/dep-tree/internal/python"
	"github.com/gabotechs/dep-tree/internal/rust"
)

const DefaultConfigPath = ".dep-tree.yml"

//go:embed sample-config.yml
var SampleConfig string

type Config struct {
	Path          string
	Exclude       []string      `yaml:"exclude"`
	UnwrapExports bool          `yaml:"unwrapExports"`
	Check         check.Config  `yaml:"check"`
	Js            js.Config     `yaml:"js"`
	Rust          rust.Config   `yaml:"rust"`
	Python        python.Config `yaml:"python"`
}

func (c *Config) UnwrapProxyExports() bool {
	return c.UnwrapExports
}

func (c *Config) IgnoreFiles() []string {
	return c.Exclude
}

func ParseConfig(cfgPath string) (*Config, error) {
	// Default values.
	cfg := Config{
		UnwrapExports: false,
		Js: js.Config{
			Workspaces:    true,
			TsConfigPaths: true,
		},
		Python: python.Config{
			ExcludeConditionalImports: false,
		},
		Rust: rust.Config{},
	}

	isDefault := cfgPath == ""
	if cfgPath == "" {
		cfgPath = DefaultConfigPath
	}
	content, err := os.ReadFile(cfgPath)
	if os.IsNotExist(err) {
		if !isDefault {
			return &cfg, err
		} else {
			return &cfg, nil
		}
	} else if err != nil {
		return &cfg, err
	}
	absCfgPath, err := filepath.Abs(cfgPath)
	if err != nil {
		return &cfg, err
	}
	cfg.Path = path.Dir(absCfgPath)

	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)
	err = decoder.Decode(&cfg)
	if err != nil {
		return &cfg, fmt.Errorf(`config file "%s" is not a valid yml file: %w`, cfgPath, err)
	}
	cfg.Check.Init(path.Dir(absCfgPath))
	return &cfg, nil
}
