//nolint:govet
package dart_grammar

import (
	"bytes"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/utils"
)

type Statement struct {
	// directives.
	Import *Import `  @@`
	Export *Export `| @@`
	PartOf *PartOf `| @@`
	Part   *Part   `| @@`
	// declarations.
	Declaration *Declaration `| @@`
}

type File struct {
	Statements []*Statement `(@@ | ANY | ALL | Punct | Ident | String | RawString | MultilineString)*`
}

var (
	lex = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"Comment", `//.*|/\*(?:.|\n)*?\*/`},
			{"MultilineString", `r?'''(?:.|\n)*?'''` + "|" + `r?"""(?:.|\n)*?"""`},
			{"RawString", `r'(?:\\.|[^'\n])*'` + "|" + `r"(?:\\.|[^"\n])*"`},
			{"String", `'(?:\\.|[^'\n])*'` + "|" + `"(?:\\.|[^"\n])*"`},
			{"ALL", `\*`},
			{"Ident", `[_$a-zA-Z][_$a-zA-Z0-9]*`},
			{"Punct", "[;:,.<>{}()\\[\\]=&|!?+\\-/%~^@]"},
			{"Whitespace", `\s+`},
			{"ANY", `.`},
		},
	)
	parser = participle.MustBuild[File](
		participle.Lexer(lex),
		participle.Elide("Whitespace", "Comment"),
		utils.UnquoteSafe("String"),
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
		AbsPath: filePath,
	}, nil
}
