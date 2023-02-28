//nolint:govet
package rust_grammar

import (
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Statement struct {
	Mod *Mod `@@`
	Use *Use `| @@`
	Pub *Pub `| @@`
}

type File struct {
	Statements []*Statement `(@@ | ANY | ALL | Punct | PathSep | Ident)*`
	Path       string
}

var (
	lex = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"ALL", `\*`},
			{"PathSep", `::`},
			{"Punct", `[,{}()]`},
			{"Ident", `[_$a-zA-Z][_$a-zA-Z0-9]*`},
			{"Comment", `//.*|/\*(.|\n)*?\*/`},
			{"Whitespace", `\s+`},
			{"ANY", `.`},
		},
	)
	parser = participle.MustBuild[File](
		participle.Lexer(lex),
		participle.Elide("Whitespace", "Comment"),
		participle.UseLookahead(1024),
	)
)

func Parse(filePath string) (*File, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	file, err := parser.ParseBytes(filePath, content)
	file.Path = filePath
	return file, err
}
