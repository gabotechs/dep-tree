//nolint:govet
package js_grammar

import (
	"os"
	"strconv"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Statement struct {
	// imports.
	DynamicImport *DynamicImport `  @@`
	StaticImport  *StaticImport  `| @@`
	// exports.
	DeclarationExport *DeclarationExport `| @@`
	DefaultExport     *DefaultExport     `| @@`
	ProxyExport       *ProxyExport       `| @@`
	ListExport        *ListExport        `| @@`
}

type File struct {
	Statements []*Statement `(@@ | ANY | ALL | Punct | Ident | String | BacktickString)*`
	Path       string
}

var (
	lex = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"ALL", `\*`},
			{"Punct", `[,{}()]`},
			{"Ident", `[_$a-zA-Z][_$a-zA-Z0-9]*`},
			{"String", `'(?:\\.|[^'])*'` + "|" + `"(?:\\.|[^"])*"`},
			{"BacktickString", "`(?:\\\\.|[^`])*`"},
			{"Comment", `//.*|/\*(.|\n)*?\*/`},
			{"Whitespace", `\s+`},
			{"ANY", `.`},
		},
	)
	parser = participle.MustBuild[File](
		participle.Lexer(lex),
		participle.Elide("Whitespace", "Comment"),
		participle.Unquote("String"),
		// Unquote BacktickStrings.
		participle.Map(func(t lexer.Token) (lexer.Token, error) {
			s := t.Value
			quote := s[0]
			s = s[1 : len(s)-1]
			out := ""
			for s != "" {
				// The `strconv.UnquoteChar` function is able to handle escaped single quotes in a single-quoted string ('/'')
				// and escaped double quotes in a double-quoted string ("/""), but it is not able to handle escaped backticks
				// in a backtick string (`/``). This conditional statement handles it.
				if len(s) >= 2 && quote == '`' && s[0] == '\\' && s[1] == '`' {
					out += string(s[1])
					s = s[2:]
					continue
				}
				value, _, tail, err := strconv.UnquoteChar(s, quote)
				if err != nil {
					return t, err
				}
				s = tail
				out += string(value)
			}
			t.Value = out
			return t, nil
		}, "BacktickString"),
		participle.UseLookahead(1024),
	)
)

func Parse(filePath string) (*File, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	file, err := parser.ParseBytes(filePath, content)
	if file != nil {
		file.Path = filePath
	}
	return file, err
}
