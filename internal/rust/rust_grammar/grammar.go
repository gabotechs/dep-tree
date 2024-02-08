//nolint:govet
package rust_grammar

import (
	"bytes"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/gabotechs/dep-tree/internal/language"
)

type Statement struct {
	Mod *Mod `@@`
	Use *Use `| @@`
	Pub *Pub `| @@`
}

type File struct {
	Statements []*Statement `(@@ | ANY | ALL | String | Punct | PathSep | Ident)*`
	Path       string
	loc        int
	size       int
}

func (f File) Loc() int {
	return f.loc
}

func (f File) Size() int {
	return f.size
}

var (
	lex = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"ALL", `\*`},
			{"PathSep", `::`},
			{"Punct", `[,{}()]`},
			{"String", `"(?:\\.|[^"])*"`},
			{"Ident", `(r#)?[_$a-zA-Z][_$a-zA-Z0-9]*`},
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

func Parse(filePath string) (*language.FileInfo, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	file, err := parser.ParseBytes(filePath, content)
	if err != nil {
		return nil, err
	}
	return &language.FileInfo{
		Content: file,
		Loc:     bytes.Count(content, []byte("\n")),
		Size:    len(content),
		Path:    filePath,
	}, nil
}
