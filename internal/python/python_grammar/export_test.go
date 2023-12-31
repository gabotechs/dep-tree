package python_grammar

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExport(t *testing.T) {
	tests := []struct {
		Name                    string
		ExpectedVariableUnpacks []VariableUnpack
		ExpectedVariableAssigns []VariableAssign
		ExpectedVariableTypings []VariableTyping
		ExpectedClasses         []Class
		ExpectedFunctions       []Function
	}{
		{
			Name:                    "foo = 1",
			ExpectedVariableAssigns: []VariableAssign{{Names: []string{"foo"}}},
		},
		{
			Name:                    "foo=1",
			ExpectedVariableAssigns: []VariableAssign{{Names: []string{"foo"}}},
		},
		{
			Name: "foo",
		},
		{
			Name: " foo = 1",
		},
		{
			Name:                    "foo: int = 1",
			ExpectedVariableTypings: []VariableTyping{{Name: "foo"}},
		},
		{
			Name:                    "foo: int",
			ExpectedVariableTypings: []VariableTyping{{Name: "foo"}},
		},
		{
			Name:                    "foo = bar = 1",
			ExpectedVariableAssigns: []VariableAssign{{Names: []string{"foo", "bar"}}},
		},
		{
			Name: " foo = bar = 1",
		},
		{
			Name:                    "foo, bar = 1, 1",
			ExpectedVariableUnpacks: []VariableUnpack{{Names: []string{"foo", "bar"}}},
		},
		{
			Name:                    "(   foo,  bar) = 1, 1",
			ExpectedVariableUnpacks: []VariableUnpack{{Names: []string{"foo", "bar"}}},
		},
		{
			Name:                    "(\n  foo,\n  bar\n) = 1, 1",
			ExpectedVariableUnpacks: []VariableUnpack{{Names: []string{"foo", "bar"}}},
		},
		{
			Name: " (\n  foo,\n  bar\n) = 1, 1",
		},
		{
			Name:              "def func():",
			ExpectedFunctions: []Function{{Name: "func"}},
		},
		{
			Name:              "async def func():",
			ExpectedFunctions: []Function{{Name: "func"}},
		},
		{
			Name: " def func():",
		},
		{
			Name:            "class cls:",
			ExpectedClasses: []Class{{Name: "cls"}},
		},
		{
			Name:            "class cls:\n    \"\"\" comment \"\"\"\n    pass",
			ExpectedClasses: []Class{{Name: "cls"}},
		},
		{
			Name: " class cls:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			var content []byte
			if strings.HasSuffix(tt.Name, ".py") {
				var err error
				content, err = os.ReadFile(filepath.Join(".export_test", tt.Name))
				a.NoError(err)
			} else {
				content = []byte(tt.Name)
			}
			parsed, err := parser.ParseBytes("", content)
			a.NoError(err)

			var variables []VariableAssign
			var unpacks []VariableUnpack
			var typings []VariableTyping
			var classes []Class
			var functions []Function
			for _, stmt := range parsed.Statements {
				switch {
				case stmt.VariableAssign != nil:
					variables = append(variables, *stmt.VariableAssign)
				case stmt.VariableUnpack != nil:
					unpacks = append(unpacks, *stmt.VariableUnpack)
				case stmt.VariableTyping != nil:
					typings = append(typings, *stmt.VariableTyping)
				case stmt.Function != nil:
					functions = append(functions, *stmt.Function)
				case stmt.Class != nil:
					classes = append(classes, *stmt.Class)
				}
			}

			a.Equal(tt.ExpectedVariableAssigns, variables)
			a.Equal(tt.ExpectedVariableUnpacks, unpacks)
			a.Equal(tt.ExpectedVariableTypings, typings)
			a.Equal(tt.ExpectedClasses, classes)
			a.Equal(tt.ExpectedFunctions, functions)
		})
	}
}
