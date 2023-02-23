package config

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Path                      string
	Entrypoints               []string            `yaml:"entrypoints"`
	AllowCircularDependencies bool                `yaml:"allowCircularDependencies"`
	Aliases                   map[string][]string `yaml:"aliases"`
	WhiteList                 map[string][]string `yaml:"allow"`
	BlackList                 map[string][]string `yaml:"deny"`
}

func ParseConfig(cfgPath string) (*Config, error) {
	if cfgPath == "" {
		cfgPath = ".dep-tree.yml"
	}

	content, err := os.ReadFile(cfgPath)
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
		return nil, err
	} else if len(cfg.Entrypoints) == 0 {
		return nil, errors.New("config file has no entrypoints")
	}
	cfg.expandAliases()
	return &cfg, nil
}

func (c *Config) expandAliases() {
	lists := []map[string][]string{
		c.WhiteList,
		c.BlackList,
	}
	for _, list := range lists {
		for k, v := range list {
			newV := make([]string, 0)
			for _, entry := range v {
				if alias, ok := c.Aliases[entry]; ok {
					newV = append(newV, alias...)
				} else {
					newV = append(newV, entry)
				}
			}
			list[k] = newV
		}
	}
}

func (c *Config) whiteListCheck(from, to string) (bool, error) {
	for k, v := range c.WhiteList {
		doesMatch, err := match(k, from)
		if err != nil {
			return false, err
		}
		if doesMatch {
			for _, dest := range v {
				shouldPass, err := match(dest, to)
				if err != nil {
					return false, err
				}
				if shouldPass {
					return true, nil
				}
			}
			return false, nil
		}
	}
	return true, nil
}

func (c *Config) blackListCheck(from, to string) (bool, error) {
	for k, v := range c.BlackList {
		doesMatch, err := match(k, from)
		if err != nil {
			return false, err
		}
		if doesMatch {
			for _, dest := range v {
				shouldReject, err := match(dest, to)
				if err != nil {
					return false, err
				}
				if shouldReject {
					return false, nil
				}
			}
		}
	}

	return true, nil
}

func (c *Config) Check(from, to string) (bool, error) {
	pass, err := c.blackListCheck(from, to)
	if err != nil || !pass {
		return pass, err
	}
	return c.whiteListCheck(from, to)
}

func (c *Config) rel(p string) string {
	relPath, err := filepath.Rel(c.Path, p)
	if err != nil {
		return p
	}
	return relPath
}

func (c *Config) validate(
	start string,
	destinations func(from string) []string,
	seen map[string]bool,
) ([]string, error) {
	collectedErrors := make([]string, 0)

	if _, ok := seen[start]; ok {
		return collectedErrors, nil
	} else {
		seen[start] = true
	}

	for _, dest := range destinations(start) {
		from, to := c.rel(start), c.rel(dest)
		pass, err := c.Check(from, to)
		if err != nil {
			return nil, err
		} else if !pass {
			collectedErrors = append(collectedErrors, from+" -> "+to)
		}
		moreErrors, err := c.validate(dest, destinations, seen)
		if err != nil {
			return nil, err
		}
		collectedErrors = append(collectedErrors, moreErrors...)
	}
	return collectedErrors, nil
}

func (c *Config) Validate(start string, destinations func(from string) []string) ([]string, error) {
	return c.validate(start, destinations, map[string]bool{})
}
