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
	FollowReExports bool          `yaml:"followReExports"`
	Check           check.Config  `yaml:"check"`
	Js              js.Config     `yaml:"js"`
	Rust            rust.Config   `yaml:"rust"`
	Python          python.Config `yaml:"python"`
}

func (c *Config) UnwrapProxyExports() bool {
	return c.FollowReExports
}

func (c *Config) IgnoreFiles() []string {
	return c.Exclude
}

func ParseConfig(cfgPath string) (*Config, error) {
	// Default values.
	cfg := Config{
		FollowReExports: false,
		Js: js.Config{
			FollowWorkspaces:    true,
			FollowTsConfigPaths: true,
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
		}
	} else if err != nil {
		return &cfg, err
	}
	absCfgPath, err := filepath.Abs(cfgPath)
	if err != nil {
		return &cfg, err
	}
	cfg.Path = path.Dir(absCfgPath)

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return &cfg, fmt.Errorf(`config file "%s" is not a valid yml file`, cfgPath)
	}
	cfg.Check.Init(path.Dir(absCfgPath))
	return &cfg, nil
}
