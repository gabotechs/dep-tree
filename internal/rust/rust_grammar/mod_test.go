package rust_grammar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMod(t *testing.T) {
	tests := []struct {
		Name         string
		ExpectedMods []Mod
	}{
		{
			Name: "mod mod",
			ExpectedMods: []Mod{{
				Name: "mod",
			}},
		},
		{
			Name: "pub mod mod",
			ExpectedMods: []Mod{{
				Pub:  true,
				Name: "mod",
			}},
		},
		{
			Name: "pub(in crate::this) mod mod",
			ExpectedMods: []Mod{{
				Pub:  true,
				Name: "mod",
			}},
		},
		{
			Name: "mod mod {}",
			ExpectedMods: []Mod{{
				Name:  "mod",
				Local: true,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			content := []byte(tt.Name)
			parsed, err := parser.ParseBytes("", content)
			a.NoError(err)

			var mods []Mod
			for _, stmt := range parsed.Statements {
				if stmt.Mod != nil {
					mods = append(mods, *stmt.Mod)
				}
			}

			a.Equal(tt.ExpectedMods, mods)
		})
	}
}
