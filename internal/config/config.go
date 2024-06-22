package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/utils"
	"gopkg.in/yaml.v3"

	"github.com/gabotechs/dep-tree/internal/check"
	golang "github.com/gabotechs/dep-tree/internal/go"
	"github.com/gabotechs/dep-tree/internal/js"
	"github.com/gabotechs/dep-tree/internal/python"
	"github.com/gabotechs/dep-tree/internal/rust"
)

const DefaultConfigPath = ".dep-tree.yml"

//go:embed sample-config.yml
var SampleConfig string

type Config struct {
	Path          string
	Source        string
	Exclude       []string      `yaml:"exclude"`
	Only          []string      `yaml:"only"`
	UnwrapExports bool          `yaml:"unwrapExports"`
	Check         check.Config  `yaml:"check"`
	Js            js.Config     `yaml:"js"`
	Rust          rust.Config   `yaml:"rust"`
	Python        python.Config `yaml:"python"`
	Golang        golang.Config `yaml:"golang"`
}

func NewConfigCwd() Config {
	var cwd, _ = os.Getwd()
	return Config{Path: cwd}
}

func (c *Config) EnsureAbsPaths() {
	for i, file := range c.Exclude {
		if !filepath.IsAbs(file) {
			c.Exclude[i] = filepath.Join(c.Path, file)
		}
	}

	for i, file := range c.Only {
		if !filepath.IsAbs(file) {
			c.Only[i] = filepath.Join(c.Path, file)
		}
	}
}

func (c *Config) ValidatePatterns() error {
	for _, pattern := range c.Exclude {
		if _, err := utils.GlobstarMatch(pattern, ""); err != nil {
			return fmt.Errorf("exclude pattern '%s' is not correctly formatted", pattern)
		}
	}

	for _, pattern := range c.Only {
		if _, err := utils.GlobstarMatch(pattern, ""); err != nil {
			return fmt.Errorf("only pattern '%s' is not correctly formatted", pattern)
		}
	}

	return nil
}

func ParseConfigFromFile(cfgPath string) (*Config, error) {
	// Default values.
	cfg := Config{
		Path:          cfgPath,
		Source:        "default",
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

	var isDefault bool
	if cfgPath == "" {
		isDefault = true
		cfgPath = DefaultConfigPath
	}
	absCfgPath, err := filepath.Abs(cfgPath)
	if err != nil {
		return nil, err
	}
	cfg.Path = filepath.Dir(absCfgPath)
	cfg.Check.Path = cfg.Path

	// If a specific path was requested, and it does not exist, fail
	// If no specific path was requested, and the default config path does not exist, succeed
	content, err := os.ReadFile(cfgPath)
	if os.IsNotExist(err) {
		if isDefault {
			return &cfg, nil
		}
		return nil, err
	} else if err != nil {
		return nil, err
	}
	cfg.Source = "file"

	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf(`config file "%s" is not a valid yml file: %w`, cfgPath, err)
	}

	cfg.Check.Init(filepath.Dir(absCfgPath))
	return &cfg, nil
}
