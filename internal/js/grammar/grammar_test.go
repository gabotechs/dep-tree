package grammar

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGrammar(t *testing.T) {
	tests := []struct {
		Name            string
		Content         string
		ExpectedStatic  []string
		ExpectedDynamic []string
	}{
		{
			Name:           "import * from 'file'",
			ExpectedStatic: []string{"file"},
		},
		{
			Name:           "import * as Something from \"file\"",
			ExpectedStatic: []string{"file"},
		},
		{
			Name:           "import * as Something from \"file\";",
			ExpectedStatic: []string{"file"},
		},
		{
			Name:           "import Something from \"file\"",
			ExpectedStatic: []string{"file"},
		},
		{
			Name:           "import { Something } from 'file'",
			ExpectedStatic: []string{"file"},
		},
		{
			Name:           "import { One, Other } from 'file'",
			ExpectedStatic: []string{"file"},
		},
		{
			Name:           "import Default, { One, Other } from 'file'",
			ExpectedStatic: []string{"file"},
		},
		{
			Name:           "import \"something\"",
			ExpectedStatic: []string{"something"},
		},
		{
			Name:            "import('something')",
			ExpectedDynamic: []string{"something"},
		},
		{
			Name:            "(import('something'))",
			ExpectedDynamic: []string{"something"},
		},
		{
			Name:           "import { One } from 'one'\nimport \"two\"",
			ExpectedStatic: []string{"one", "two"},
		},
		{
			Name:    "All imports",
			Content: "import-regex.js",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			var content []byte
			if strings.HasSuffix(tt.Content, ".js") {
				var err error
				content, err = os.ReadFile(path.Join("grammar_test", tt.Content))
				a.NoError(err)
			} else if tt.Content != "" {
				content = []byte(tt.Content)
			} else {
				content = []byte(tt.Name)
			}
			parsed, err := Parser.ParseBytes("", content)
			a.NoError(err)

			var staticResults []string
			var dynamicResults []string
			for _, imp := range parsed.Imports {
				if imp.StaticImport != nil {
					staticResults = append(staticResults, imp.StaticImport.Path)
				} else if imp.DynamicImport != nil {
					dynamicResults = append(dynamicResults, imp.DynamicImport.Path)
				}
			}
			a.Equal(tt.ExpectedStatic, staticResults)
			a.Equal(tt.ExpectedDynamic, dynamicResults)
		})
	}
}
