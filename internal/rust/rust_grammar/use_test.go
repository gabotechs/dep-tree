package rust_grammar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUse(t *testing.T) {
	tests := []struct {
		Name        string
		ExpectedUse []Use
	}{
		{
			Name: "use something::something;",
			ExpectedUse: []Use{{
				PathSlices: []string{"something"},
				Names:      []Name{{Original: "something"}},
			}},
		},
		{
			Name: "use something::something as another;",
			ExpectedUse: []Use{{
				PathSlices: []string{"something"},
				Names:      []Name{{Original: "something", Alias: "another"}},
			}},
		},
		{
			Name: "pub use something::something;",
			ExpectedUse: []Use{{
				Pub:        true,
				PathSlices: []string{"something"},
				Names:      []Name{{Original: "something"}},
			}},
		},
		{
			Name: "pub (   crate) use something  ::something  ;",
			ExpectedUse: []Use{{
				Pub:        true,
				PathSlices: []string{"something"},
				Names:      []Name{{Original: "something"}},
			}},
		},
		{
			Name: "use something::{something};",
			ExpectedUse: []Use{{
				PathSlices: []string{"something"},
				Names:      []Name{{Original: "something"}},
			}},
		},
		{
			Name: "use something::{one, OrAnother};",
			ExpectedUse: []Use{{
				PathSlices: []string{"something"},
				Names:      []Name{{Original: "one"}, {Original: "OrAnother"}},
			}},
		},
		{
			Name: "use something::{one as two, OrAnother};",
			ExpectedUse: []Use{{
				PathSlices: []string{"something"},
				Names:      []Name{{Original: "one", Alias: "two"}, {Original: "OrAnother"}},
			}},
		},
		{
			Name: "use one::very_long::veeery_long::path::something;",
			ExpectedUse: []Use{{
				PathSlices: []string{"one", "very_long", "veeery_long", "path"},
				Names:      []Name{{Original: "something"}},
			}},
		},
		{
			Name: "use one::very_long::veeery_long::path::*;",
			ExpectedUse: []Use{{
				PathSlices: []string{"one", "very_long", "veeery_long", "path"},
				All:        true,
			}},
		},
		{
			Name: "pub use crate::ast::{Node, Operator};\n",
			ExpectedUse: []Use{{
				Pub:        true,
				PathSlices: []string{"crate", "ast"},
				Names:      []Name{{Original: "Node"}, {Original: "Operator"}},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			content := []byte(tt.Name)
			parsed, err := parser.ParseBytes("", content)
			a.NoError(err)

			var uses []Use
			for _, stmt := range parsed.Statements {
				if stmt.Use != nil {
					uses = append(uses, *stmt.Use)
				}
			}

			a.Equal(tt.ExpectedUse, uses)
		})
	}
}
