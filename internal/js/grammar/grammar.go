//nolint:govet
package grammar

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Statement struct {
	// imports
	DynamicImport *DynamicImport `  @@`
	StaticImport  *StaticImport  `| @@`
	// exports
	DeclarationExport *DeclarationExport `| @@`
	DefaultExport     *DefaultExport     `| @@`
	ProxyExport       *ProxyExport       `| @@`
	ListExport        *ListExport        `| @@`
}

type File struct {
	Statements []*Statement `(@@ | ANY | ALL | Punct | Ident | String)*`
}

var (
	lex = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"ALL", `\*`},
			{"Punct", `[,{}()]`},
			{"Ident", `[_$a-zA-Z\\xA0-\\uFFFF][_$a-zA-Z0-9\\xA0-\\uFFFF]*`},
			{"String", `'[^']*'|"[^"]*"`},
			{"Comment", `//.*|/\*.*?\*/`},
			{"Whitespace", `\s+`},
			{"ANY", `.`},
		},
	)
	parser = participle.MustBuild[File](
		participle.Lexer(lex),
		participle.Elide("Whitespace", "Comment"),
		participle.Unquote("String"),
		participle.UseLookahead(1024),
	)
)

func Parse(content []byte) (*File, error) {
	return parser.ParseBytes("", content)
}
