package python_grammar

import (
	"github.com/stretchr/testify/require"

	"os"
	"path"
	"strings"
	"testing"
)

func TestExport(t *testing.T) {
	tests := []struct {
		Name              string
		ExpectedVariables []Variable
		ExpectedClasses   []Class
		ExpectedFunctions []Function
	}{
		{
			Name:              "foo = 1",
			ExpectedVariables: []Variable{{Name: "foo"}},
		},
		{
			Name:              "foo=1",
			ExpectedVariables: []Variable{{Name: "foo"}},
		},
		{
			Name:              "foo",
			ExpectedVariables: nil,
		},
		{
			Name:              " foo = 1",
			ExpectedVariables: []Variable{{Name: "foo", Indented: true}},
		},
		{
			Name:              "foo: int = 1",
			ExpectedVariables: []Variable{{Name: "foo"}},
		},
		{
			Name:              "foo: int",
			ExpectedVariables: []Variable{{Name: "foo"}},
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
			Name:              " def func():",
			ExpectedFunctions: []Function{{Name: "func", Indented: true}},
		},
		{
			Name:            "class cls:",
			ExpectedClasses: []Class{{Name: "cls"}},
		},
		{
			Name:            " class cls:",
			ExpectedClasses: []Class{{Name: "cls", Indented: true}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			var content []byte
			if strings.HasSuffix(tt.Name, ".py") {
				var err error
				content, err = os.ReadFile(path.Join(".export_test", tt.Name))
				a.NoError(err)
			} else {
				content = []byte(tt.Name)
			}
			parsed, err := parser.ParseBytes("", content)
			a.NoError(err)

			var variables []Variable
			var classes []Class
			var functions []Function
			for _, stmt := range parsed.Statements {
				switch {
				case stmt.Variable != nil:
					variables = append(variables, *stmt.Variable)
				case stmt.Function != nil:
					functions = append(functions, *stmt.Function)
				case stmt.Class != nil:
					classes = append(classes, *stmt.Class)
				}
			}
			a.Equal(tt.ExpectedVariables, variables)
			a.Equal(tt.ExpectedClasses, classes)
			a.Equal(tt.ExpectedFunctions, functions)
		})
	}
}
