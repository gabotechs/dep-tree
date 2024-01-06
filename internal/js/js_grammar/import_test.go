package js_grammar

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImport(t *testing.T) {
	tests := []struct {
		Name            string
		ExpectedStatic  []string
		ExpectedDynamic []string
		ExpectedRequire []string
	}{
		{
			Name:           "import * from 'file'",
			ExpectedStatic: []string{"file"},
		},
		{
			Name:           "// comment \nimport * from 'file'",
			ExpectedStatic: []string{"file"},
		},
		{
			Name:           "useless shit, import * from 'file' dumb suffix",
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
			Name:           "import type { Something } from 'file'",
			ExpectedStatic: []string{"file"},
		},
		{
			Name:           "import { One, Other } from 'file'",
			ExpectedStatic: []string{"file"},
		},
		{
			Name:           "import { One, Other, } from 'file'",
			ExpectedStatic: []string{"file"},
		},
		{
			Name:           "import { type One, Other } from 'file'",
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
			Name:            "import   ('something'); const a",
			ExpectedDynamic: []string{"something"},
		},
		{
			Name:           "import { One } from 'one'\nimport \"two\"",
			ExpectedStatic: []string{"one", "two"},
		},
		{
			Name:           "import { from } from 'somewhere'",
			ExpectedStatic: []string{"somewhere"},
		},
		{
			Name:           "import { from, } from 'somewhere'",
			ExpectedStatic: []string{"somewhere"},
		},
		{
			Name:           "const importVariable = []\nimport 'variable'",
			ExpectedStatic: []string{"variable"},
		},
		{
			Name:           "import '.export'",
			ExpectedStatic: []string{".export"},
		},
		{
			Name:            "require('foo')",
			ExpectedRequire: []string{"foo"},
		},
		{
			Name:            "const a = require('foo')",
			ExpectedRequire: []string{"foo"},
		},
		{
			Name:            "var a = require('foo');",
			ExpectedRequire: []string{"foo"},
		},
		{
			Name:            "const { a, b } = require('foo')",
			ExpectedRequire: []string{"foo"},
		},
		{
			Name:            "const { a: a1, b } = require('foo')",
			ExpectedRequire: []string{"foo"},
		},
		{
			Name:            "const { a: a1, b: b2, } = require('foo')",
			ExpectedRequire: []string{"foo"},
		},
		{
			Name:            "let { a, b, } = require('foo')",
			ExpectedRequire: []string{"foo"},
		},
		{
			Name: "import-regex.js",
			ExpectedStatic: []string{
				"@angular2/core",
				"module-name",
				"module-name  ",
				"  module-name",
				"module-name",
				"module-name",
				"module-name",
				"@angular2/core",
				"$module-name",
				"module-name",
				"module-name",
				"module-name  ",
				"  module-name",
				"module-name",
				"module-name",
				"module-name",
				"@angular2/core",
				"$module-name",
				"module-name",
				"module-name",
				"module-name",
				"react",
				"redux-form",
				"module-name",
				"../geometries/Geometries.js",
				"../geometries/Geometries.js",
				"redux-form",
				"./views/ListView",
				"./views/AddView",
				"./views/EditView",
				"redux-form",
				"./views/ListView",
				"./views/AddView",
				"./views/EditView",
			},
			ExpectedDynamic: []string{
				"whatever.js",
				"whatever.js",
			},
		},
		{
			Name: "test-1.js",
			ExpectedStatic: []string{
				"react",
				"../services/apollo",
				"../config",
				"../styles/theme",
				"history",
				"../constants/routing",
				"../components/dialogs/SnackBar",
				"@apollo/client",
				"react-router",
				"./contexts/AppContext",
				"@material-ui/core",
				"../views/SlicerView/contexts/SlicingContext",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			var content []byte
			if strings.HasSuffix(tt.Name, ".js") {
				var err error
				content, err = os.ReadFile(filepath.Join(".import_test", tt.Name))
				a.NoError(err)
			} else {
				content = []byte(tt.Name)
			}
			parsed, err := parser.ParseBytes("", content)
			a.NoError(err)

			var staticResults []string
			var dynamicResults []string
			var requires []string
			for _, stmt := range parsed.Statements {
				switch {
				case stmt.StaticImport != nil:
					staticResults = append(staticResults, stmt.StaticImport.Path)
				case stmt.DynamicImport != nil:
					dynamicResults = append(dynamicResults, stmt.DynamicImport.Path)
				case stmt.Require != nil:
					requires = append(requires, stmt.Require.Path)
				}
			}
			a.Equal(tt.ExpectedStatic, staticResults)
			a.Equal(tt.ExpectedDynamic, dynamicResults)
			a.Equal(tt.ExpectedRequire, requires)
		})
	}
}
