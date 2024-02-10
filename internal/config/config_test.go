package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testFolder = ".config_test"

func TestParseConfig(t *testing.T) {
	absTestFolder, _ := filepath.Abs(testFolder)

	tests := []struct {
		Name              string
		File              string
		ExpectedWhiteList map[string][]string
		ExpectedBlackList map[string][]string
		ExpectedExclude   []string
	}{
		{
			Name: "default file",
			File: "",
		},
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
		{
			Name: "Exclusion",
			File: ".excludes.yml",
			ExpectedExclude: []string{
				filepath.Join(absTestFolder, "**/foo.js"),
				filepath.Join(absTestFolder, "*/foo.js"),
				filepath.Join(absTestFolder, "foo/**/foo.js"),
				"/**/*.js",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			if tt.File != "" {
				tt.File = filepath.Join(testFolder, tt.File)
			}
			cfg, err := ParseConfig(tt.File)
			a.NoError(err)

			a.Equal(tt.ExpectedWhiteList, cfg.Check.WhiteList)
			a.Equal(tt.ExpectedBlackList, cfg.Check.BlackList)
			a.Equal(tt.ExpectedExclude, cfg.Exclude)
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
			File:     ".non-existing.yml",
			Expected: "no such file or directory",
		},
		{
			Name:     "Invalid yml",
			File:     filepath.Join(testFolder, ".invalid.yml"),
			Expected: "not a valid yml file",
		},
		{
			Name:     "Entrypoints on top level",
			File:     filepath.Join(testFolder, ".top-level-entrypoints.yml"),
			Expected: "not a valid yml file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_, err := ParseConfig(tt.File)
			a.ErrorContains(err, tt.Expected)
		})
	}
}

func TestSampleConfig(t *testing.T) {
	a := require.New(t)
	_, err := ParseConfig("sample-config.yml")
	a.NoError(err)
}
