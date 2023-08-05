package js_grammar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGrammar(t *testing.T) {
	tests := []struct {
		Name    string
		Content string
	}{
		{
			Name:    "Double quoted string with inner double quotes",
			Content: `export const a = "this is a \"string\""`,
		},
		{
			Name:    "Quoted string with inner quotes",
			Content: `export const a = 'this is a \'string\''`,
		},
		{
			Name:    "Quoted string with inner escaped quotes",
			Content: `'\\\\').replace(/'/g, "\\'")`,
		},
		{
			Name:    "Doc with quotes",
			Content: "/* EscapedCharacter :: one of `\"` `\\` `/` `b` `f` `n` `r` `t` \n*/\nimport { MuiThemeProvider } from \"@material-ui/core\";",
		},
		{
			Name:    "More quoted strings",
			Content: "`\"\"\" \\` \"\"\"`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			_, err := parser.ParseBytes("", []byte(tt.Content))
			a.NoError(err)
		})
	}
}
