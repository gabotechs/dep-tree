package config

import (
	"path/filepath"
	"testing"

	"github.com/gabotechs/dep-tree/internal/check"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testFolder = ".config_test"

func TestParseConfig(t *testing.T) {
	absTestFolder, _ := filepath.Abs(testFolder)

	tests := []struct {
		Name              string
		File              string
		ExpectedWhiteList map[string]check.WhiteListEntries
		ExpectedBlackList map[string][]check.BlackListEntry
		ExpectedExclude   []string
	}{
		{
			Name: "default file",
			File: "",
		},
		{
			Name: "Simple",
			File: ".parse.yml",
			ExpectedWhiteList: map[string]check.WhiteListEntries{
				"foo": {To: []string{"bar"}},
			},
			ExpectedBlackList: map[string][]check.BlackListEntry{
				"bar": {{To: "baz"}},
			},
		},
		{
			Name: "Aliased",
			File: ".aliases.yml",
			ExpectedWhiteList: map[string]check.WhiteListEntries{
				"src/users/**": {To: []string{
					"src/users/**",
					"src/@*/**",
					"src/config.js",
					"src/common.js",
				}},
			},
			ExpectedBlackList: map[string][]check.BlackListEntry{
				"src/products/**": {
					{To: "src/users/**"},
					{To: "src/@*/**"},
					{To: "src/config.js"},
					{To: "src/common.js"},
					{To: "src/generated/**"},
					{To: "generated/**"},
				},
			},
		},
		{
			Name: "With descriptions",
			File: ".with-descriptions.yml",
			ExpectedWhiteList: map[string]check.WhiteListEntries{
				"src/users/**": {
					To: []string{
						"src/users/**",
						"1.js",
						"2.js",
						"3.js",
					},
					Reason: "only users and common allowed",
				},
			},
			ExpectedBlackList: map[string][]check.BlackListEntry{
				"src/products/**": {
					{To: "src/users/**", Reason: "users not allowed"},
					{To: "src/products/**"},
					{To: "1.js", Reason: "common 1-3 not allowed,\ndouble check your dependencies\n"},
					{To: "2.js", Reason: "common 1-3 not allowed,\ndouble check your dependencies\n"},
					{To: "3.js", Reason: "common 1-3 not allowed,\ndouble check your dependencies\n"},
					{To: "4.js"},
					{To: "5.js"},
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
			if tt.File != "" {
				tt.File = filepath.Join(testFolder, tt.File)
			}
			cfg, err := ParseConfigFromFile(tt.File)
			require.NoError(t, err)
			cfg.EnsureAbsPaths()

			assert.Equal(t, tt.ExpectedWhiteList, cfg.Check.WhiteList)
			assert.Equal(t, tt.ExpectedBlackList, cfg.Check.BlackList)
			assert.Equal(t, tt.ExpectedExclude, cfg.Exclude)
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
			_, err := ParseConfigFromFile(tt.File)
			a.ErrorContains(err, tt.Expected)
		})
	}
}

func TestSampleConfig(t *testing.T) {
	a := require.New(t)
	_, err := ParseConfigFromFile("sample-config.yml")
	a.NoError(err)
}
