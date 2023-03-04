package config

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const testFolder = ".config_test"

func TestParseConfig(t *testing.T) {
	tests := []struct {
		Name              string
		File              string
		ExpectedWhiteList map[string][]string
		ExpectedBlackList map[string][]string
	}{
		{
			Name: "Simple",
			File: ".parse.yml",
			ExpectedWhiteList: map[string][]string{
				"foo": {"bar"},
			},
			ExpectedBlackList: map[string][]string{
				"bar": {"baz"},
			},
		},
		{
			Name: "Aliased",
			File: ".aliases.yml",
			ExpectedWhiteList: map[string][]string{
				"src/users/**": {
					"src/users/**",
					"src/@*/**",
					"src/config.js",
					"src/common.js",
				},
			},
			ExpectedBlackList: map[string][]string{
				"src/products/**": {
					"src/users/**",
					"src/@*/**",
					"src/config.js",
					"src/common.js",
					"src/generated/**",
					"generated/**",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			cfg, err := ParseConfig(path.Join(testFolder, tt.File))
			a.NoError(err)

			a.Equal(tt.ExpectedWhiteList, cfg.WhiteList)
			a.Equal(tt.ExpectedBlackList, cfg.BlackList)
		})
	}
}

func TestConfig_ErrorHandling(t *testing.T) {
	tests := []struct {
		Name     string
		File     string
		Expected string
	}{
		{
			Name:     "No config file",
			File:     "",
			Expected: "no such file or directory",
		},
		{
			Name:     "No entrypoints",
			File:     path.Join(testFolder, ".no-entrypoints.yml"),
			Expected: "has no entrypoints",
		},
		{
			Name:     "Invalid yml",
			File:     path.Join(testFolder, ".invalid.yml"),
			Expected: "not a valid yml file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_, err := ParseConfig(tt.File)
			a.Contains(err.Error(), tt.Expected)
		})
	}
}

func TestConfig_Check(t *testing.T) {
	a := require.New(t)

	tests := []struct {
		Name   string
		Config Config
		From   string
		To     string
		Passes bool
	}{
		{
			Name: "white list passes",
			Config: Config{
				WhiteList: map[string][]string{
					"white": {"**pass**"},
				},
			},
			From:   "white",
			To:     "this is going to pass",
			Passes: true,
		},
		{
			Name: "white list fails",
			Config: Config{
				WhiteList: map[string][]string{
					"white": {"**pass**"},
				},
			},
			From:   "white",
			To:     "this doesn't",
			Passes: false,
		},
		{
			Name: "black list passes",
			Config: Config{
				BlackList: map[string][]string{
					"black": {"**fail**"},
				},
			},
			From:   "black",
			To:     "this is going to pass",
			Passes: true,
		},
		{
			Name: "black list fails",
			Config: Config{
				BlackList: map[string][]string{
					"black": {"**fail**"},
				},
			},
			From:   "black",
			To:     "this is going to fail",
			Passes: false,
		},
		{
			Name: "this should never pass",
			Config: Config{
				BlackList: map[string][]string{
					"**": {"**"},
				},
			},
			From:   "black",
			To:     "this is going to fail",
			Passes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			pass, err := tt.Config.Check(tt.From, tt.To)
			a.NoError(err)
			a.Equal(tt.Passes, pass)
		})
	}
}
