//nolint:govet
package grammar

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Deconstruction struct {
	Names []string `"{" @Ident ("as" Ident)? ("," (@Ident ("as" Ident)?)?)* "}"`
}

type AllImport struct {
	Alias string `ALL ("as" @Ident)?`
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
	Imported *Imported `"import" (@@ "from")?`
	Path     string    `@String`
}

type DynamicImport struct {
	Path string `"import" "(" @String ")"`
}
type Statement struct {
	DynamicImport *DynamicImport `@@`
	StaticImport  *StaticImport  `| @@`
}

type File struct {
	Statements []*Statement `(@@ | ANY | FALSE_IMPORT_1 | FALSE_IMPORT_2)*`
}

var (
	lex = lexer.MustStateful(lexer.Rules{
		"Root": {
			{"FALSE_IMPORT_1", `import[_$a-zA-Z0-9\\xA0-\\uFFFF]+`, nil},
			{"FALSE_IMPORT_2", `[_$a-zA-Z0-9\\xA0-\\uFFFF]+import`, nil},
			{"IMPORT", `import`, lexer.Push("Import")},
			{"Whitespace", `\s+`, nil},
			{"ANY", `.`, nil},
		},
		"Import": {
			// Keywords.
			{"COMMA", ",", nil},
			{"ALL", `\*`, nil},
			{"BRACKET_L", `{`, nil},
			{"BRACKET_R", `}`, nil},
			{"PARENTHESIS_L", `\(`, nil},
			{"PARENTHESIS_R", `\)`, nil},
			// Other.
			{"Ident", `[_$a-zA-Z\\xA0-\\uFFFF][_$a-zA-Z0-9\\xA0-\\uFFFF]*`, nil},
			{"String", `'[^']*'|"[^"]*"`, lexer.Pop()},
			{"Comment", `//.*|/\*.*?\*/`, nil},
			{"Whitespace", `\s+`, nil},
		},
	})
	parser = participle.MustBuild[File](
		participle.Lexer(lex),
		participle.Elide("Whitespace", "Comment"),
		participle.Unquote("String"),
		participle.UseLookahead(2),
	)
)

func Parse(content []byte) (*File, error) {
	return parser.ParseBytes("", content)
}
