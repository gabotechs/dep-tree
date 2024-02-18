package dummy

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type ImportStatement struct {
	Symbols []string `"import" @Ident ("," @Ident)*`
	From    string   `"from" @(Ident|Punctuation|"/")*`
}

type ExportStatement struct {
	Symbol string `"export" @Ident`
}

type Statement struct {
	Import *ImportStatement `@@ |`
	Export *ExportStatement `@@`
}

type File struct {
	Statements []Statement `@@*`
}

var (
	lex = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"KewWord", "(export|import|from)"},
			{"Punctuation", `[,\./]`},
			{"Ident", `[a-zA-Z]+`},
			{"Whitespace", `\s+`},
		},
	)
	parser = participle.MustBuild[File](
		participle.Lexer(lex),
		participle.Elide("Whitespace"),
	)
)
