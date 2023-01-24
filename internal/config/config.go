package config

import (
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Path      string
	WhiteList map[string][]string `yaml:"white_list"`
	BlackList map[string][]string `yaml:"black_list"`
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
	}

	return &cfg, nil
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
	errors := make([]string, 0)

	if _, ok := seen[start]; ok {
		return errors, nil
	} else {
		seen[start] = true
	}

	for _, dest := range destinations(start) {
		from, to := c.rel(start), c.rel(dest)
		pass, err := c.Check(from, to)
		if err != nil {
			return nil, err
		} else if !pass {
			errors = append(errors, from+" -> "+to)
		}
		moreErrors, err := c.validate(dest, destinations, seen)
		if err != nil {
			return nil, err
		}
		errors = append(errors, moreErrors...)
	}
	return errors, nil
}

func (c *Config) Validate(start string, destinations func(from string) []string) ([]string, error) {
	return c.validate(start, destinations, map[string]bool{})
}
