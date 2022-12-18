//nolint:govet
package grammar

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Deconstruction struct {
	Names []string `"{" @Ident (AS Ident)? ("," @Ident (AS Ident)?)* "}"`
}

type AllImport struct {
	Alias string `ALL (AS @Ident)?`
}

type SelectionImport struct {
	AllImport      *AllImport      `(@@?`
	Deconstruction *Deconstruction ` @@?)!`
}

type Imported struct {
	Default         bool             `(@Ident? ","?`
	SelectionImport *SelectionImport ` @@?)!`
}

type StaticImport struct {
	Imported *Imported `IMPORT (@@ FROM)?`
	Path     string    `@String`
}

type DynamicImport struct {
	Path string `IMPORT "(" @String ")"`
}
type Import struct {
	DynamicImport *DynamicImport `@@`
	StaticImport  *StaticImport  `| @@`
}

type File struct {
	Imports []*Import `((@@? ANY?)!)*`
}

var (
	lex = lexer.MustSimple([]lexer.SimpleRule{
		// Keywords.
		{"IMPORT", "import"},
		{"AS", "as"},
		{"COMMA", ","},
		{"COLON", ";"},
		{"FROM", "from"},
		{"ALL", `\*`},
		{"BRACKET_L", `{`},
		{"BRACKET_R", `}`},
		{"PARENTHESIS_L", `\(`},
		{"PARENTHESIS_R", `\)`},
		// Other.
		{"Ident", `[_$a-zA-Z\\xA0-\\uFFFF][_$a-zA-Z0-9\\xA0-\\uFFFF]*`},
		{"String", `'[^']*'|"[^"]*"`},
		{"Comment", `//.*|/\*.*?\*/`},
		{"Whitespace", `\s+`},
		// Any.
		{"ANY", `.`},
	})
	Parser = participle.MustBuild[File](
		participle.Lexer(lex),
		participle.Elide("Whitespace", "Comment"),
		participle.Unquote("String"),
		participle.UseLookahead(2),
	)
)
