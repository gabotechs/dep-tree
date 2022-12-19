//nolint:govet
package grammar

import "github.com/alecthomas/participle/v2/lexer"

type ImportDeconstruction struct {
	Names []string `"{" @Ident ("as" Ident)? ("," (@Ident ("as" Ident)?)?)* "}"`
}

type AllImport struct {
	Alias string `ALL ("as" @Ident)?`
}

type SelectionImport struct {
	AllImport      *AllImport            `(@@?`
	Deconstruction *ImportDeconstruction ` @@?)!`
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

var importLexer = lexer.Rules{
	"Import": {
		{"ALL", `\*`, nil},
		{"Punct", `[,{}()]`, nil},
		{"Ident", `[_$a-zA-Z\\xA0-\\uFFFF][_$a-zA-Z0-9\\xA0-\\uFFFF]*`, nil},
		{"String", `'[^']*'|"[^"]*"`, lexer.Pop()},
		{"Comment", `//.*|/\*.*?\*/`, nil},
		{"Whitespace", `\s+`, nil},
	},
}
