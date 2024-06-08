package golang

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExports(t *testing.T) {
	tests := []struct {
		Name     string
		Expected []string
	}{
		{
			Name:     "exports.go",
			Expected: []string{},
		},
		{
			Name:     "config.go",
			Expected: []string{"Config"},
		},
		{
			Name:     "language.go",
			Expected: []string{"Extensions", "Language", "NewLanguage"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			lang, err := NewLanguage(".", &Config{})
			a.NoError(err)
			file, err := lang.ParseFile(tt.Name)
			a.NoError(err)
			exports, err := lang.ParseExports(file)
			a.NoError(err)

			actualExports := make([]string, 0)
			for _, export := range exports.Exports {
				for _, symbol := range export.Symbols {
					actualExports = append(actualExports, symbol.Original)
				}
			}
			sort.Strings(tt.Expected)
			sort.Strings(actualExports)

			a.Equal(tt.Expected, actualExports)
		})
	}
}
