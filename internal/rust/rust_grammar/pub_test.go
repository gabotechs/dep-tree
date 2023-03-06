package rust_grammar

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPub(t *testing.T) {
	tests := []struct {
		Name        string
		File        string
		ExpectedPub []Pub
	}{
		{
			Name: "pub fn my_function",
			ExpectedPub: []Pub{{
				Name: "my_function",
			}},
		},
		{
			Name: "pub unsafe fn my_function",
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
			Name: "pub enum my_enum",
			ExpectedPub: []Pub{{
				Name: "my_enum",
			}},
		},
		{
			Name: "pub type my_type",
			ExpectedPub: []Pub{{
				Name: "my_type",
			}},
		},
		{
			Name: "pub(crate) fn my_function and a lot of shit after",
			ExpectedPub: []Pub{{
				Name: "my_function",
			}},
		},
		{
			Name: "pub async fn my_function ",
			ExpectedPub: []Pub{{
				Name: "my_function",
			}},
		},
		{
			Name: "pub static VAR",
			ExpectedPub: []Pub{{
				Name: "VAR",
			}},
		},
		{
			Name: "pub const VAR",
			ExpectedPub: []Pub{{
				Name: "VAR",
			}},
		},
		{
			Name:        "\"pub type my_type\"",
			ExpectedPub: nil,
		},
		{
			Name: "' pub struct MyStruct '",
			ExpectedPub: []Pub{{
				Name: "MyStruct",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			var content []byte
			if tt.File != "" {
				var err error
				content, err = os.ReadFile(path.Join(".test_files", tt.File))
				a.NoError(err)
			} else {
				content = []byte(tt.Name)
			}
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
