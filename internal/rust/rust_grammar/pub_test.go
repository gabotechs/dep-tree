package rust_grammar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPub(t *testing.T) {
	tests := []struct {
		Name        string
		ExpectedPub []Pub
	}{
		{
			Name: "pub fn my_function",
			ExpectedPub: []Pub{{
				Name: "my_function",
			}},
		},
		{
			Name: "pub trait my_trait",
			ExpectedPub: []Pub{{
				Name: "my_trait",
			}},
		},
		{
			Name: "pub struct my_struct",
			ExpectedPub: []Pub{{
				Name: "my_struct",
			}},
		},
		{
			Name: "pub(crate) fn my_function and a lot of shit after",
			ExpectedPub: []Pub{{
				Name: "my_function",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			content := []byte(tt.Name)
			parsed, err := parser.ParseBytes("", content)
			a.NoError(err)

			var pubs []Pub
			for _, stmt := range parsed.Statements {
				if stmt.Pub != nil {
					pubs = append(pubs, *stmt.Pub)
				}
			}

			a.Equal(tt.ExpectedPub, pubs)
		})
	}
}
