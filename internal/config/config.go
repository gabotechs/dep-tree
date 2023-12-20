package config

import (
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
	Path            string
	Exclude         []string      `yaml:"exclude"`
	FollowReExports *bool         `yaml:"followReExports,omitempty"`
	Check           check.Config  `yaml:"check"`
	Js              js.Config     `yaml:"js"`
	Rust            rust.Config   `yaml:"rust"`
	Python          python.Config `yaml:"python"`
}

func (c *Config) UnwrapProxyExports() bool {
	if c.FollowReExports == nil {
		return true
	}
	return *c.FollowReExports
}

func (c *Config) IgnoreFiles() []string {
	return c.Exclude
}

func ParseConfig(cfgPath string) (*Config, error) {
	if cfgPath == "" {
		cfgPath = DefaultConfigPath
	}

	content, err := os.ReadFile(cfgPath)
	if os.IsNotExist(err) {
		return &Config{}, err
	}
	if err != nil {
		return nil, err
	}
	absCfgPath, err := filepath.Abs(cfgPath)
	if err != nil {
		return nil, err
	}

	cfg := Config{
		Path: path.Dir(absCfgPath),
	}

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return nil, fmt.Errorf(`config file "%s" is not a valid yml file`, cfgPath)
	}
	cfg.Check.Init(path.Dir(absCfgPath))
	return &cfg, nil
}
