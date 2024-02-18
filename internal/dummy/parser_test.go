package dummy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	tests := []struct {
		Name       string
		Statements []Statement
	}{
		{
			Name: "import foo from bar",
			Statements: []Statement{{
				Import: &ImportStatement{[]string{"foo"}, "bar"},
			}},
		},
		{
			Name: "import foo, baz from ./bar.dl",
			Statements: []Statement{{
				Import: &ImportStatement{[]string{"foo", "baz"}, "./bar.dl"},
			}},
		},
		{
			Name: "export foo",
			Statements: []Statement{{
				Export: &ExportStatement{"foo"},
			}},
		},

		{
			Name: "import foo, baz from ./bar.dl",
			Statements: []Statement{{
				Import: &ImportStatement{[]string{"foo", "baz"}, "./bar.dl"},
			}},
		},
		{
			Name: "import foo, baz from ./bar.dl\n\nexport foo",
			Statements: []Statement{
				{Import: &ImportStatement{[]string{"foo", "baz"}, "./bar.dl"}},
				{Export: &ExportStatement{"foo"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			result, err := parser.ParseBytes("", []byte(tt.Name))
			a.NoError(err)

			a.Equal(tt.Statements, result.Statements)
		})
	}
}
