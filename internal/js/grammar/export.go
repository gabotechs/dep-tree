//nolint:govet
package grammar

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type ExportDeconstruction struct {
	Names []string `"{" ((Ident "as" @Ident) | @Ident) ("," ((Ident "as" @Ident) | @Ident))* "}"`
}

type DeclarationExport struct {
	Default bool   `"export" ("let"|"const"|"var"|"function"|"class")`
	Name    string `@Ident`
}

type ListExport struct {
	ExportDeconstruction *ExportDeconstruction `"export" @@`
}

type DefaultExport struct {
	Default bool `"export" "default"`
}

type ProxyExport struct {
	ExportDeconstruction *ExportDeconstruction `"export" @@`
	From                 string                `"from" @String`
}

var exportLexer = lexer.Rules{
	"CommonExport": {
		{"ALL", `\*`, nil},
		{"Comment", commentRe, nil},
		{"Whitespace", `\s+`, nil},
		{"Punct", punctuationRe, nil},
	},
	"DeclarationExport": {
		{Name: "String", Pattern: stringRe},
		{Name: "Ident", Pattern: identRe, Action: lexer.Pop()},
		lexer.Include("CommonExport"),
	},
	"ListExport": {
		{Name: "ClosingBracket", Pattern: `}`, Action: lexer.Pop()},
		{Name: "String", Pattern: stringRe},
		{Name: "Ident", Pattern: identRe},
		lexer.Include("CommonExport"),
	},
	"DefaultExport": {
		{Name: "Default", Pattern: `default`, Action: lexer.Pop()},
		{Name: "String", Pattern: stringRe},
		{Name: "Ident", Pattern: identRe},
		lexer.Include("CommonExport"),
	},
	"ProxyExport": {
		{Name: "String", Pattern: stringRe, Action: lexer.Pop()},
		{Name: "Ident", Pattern: identRe},
		lexer.Include("CommonExport"),
	},
}
