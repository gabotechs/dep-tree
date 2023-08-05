//nolint:govet
package js_grammar

import (
	"os"

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
	Statements []*Statement `(@@ | ANY | ALL | Punct | Ident | String)*`
	Path       string
}

var (
	lex = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"ALL", `\*`},
			{"Punct", `[,{}()]`},
			{"Ident", `[_$a-zA-Z][_$a-zA-Z0-9]*`},
			{"String", `'(?:\\.|[^'])*'|"(?:\\.|[^"])*"`},
			{"Comment", `//.*|/\*(.|\n)*?\*/`},
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
