//nolint:govet
package python_grammar

import (
	"bytes"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/gabotechs/dep-tree/internal/language"
)

type Statement struct {
	// imports.
	FromImport *FromImport `@@ |`
	Import     *Import     `@@ |`
	// exports.
	Class          *Class          `@@ |`
	Function       *Function       `@@ |`
	VariableTyping *VariableTyping `@@ |`
	VariableUnpack *VariableUnpack `@@ |`
	VariableAssign *VariableAssign `@@`
}

type File struct {
	Statements []*Statement `(@@ | SOF | ANY | ALL | Ident | Space | NewLine | String | MultilineString)*`
}

var (
	lex = LexerWithSOF(lexer.MustSimple(
		[]lexer.SimpleRule{
			{"ALL", `\*`},
			{"Ident", `[_$a-zA-Z][_$a-zA-Z0-9]*`},
			{"MultilineString", `'''(.|\n)*?'''` + "|" + `"""(.|\n)*?"""`},
			{"String", `'(?:\\.|[^'])*'` + "|" + `"(?:\\.|[^"])*"`},
			// https://stackoverflow.com/questions/69184441/regular-expression-for-comments-in-python-re
			{"Comment", `#.*`},
			{"NewLine", `[\n\r]`},
			{"Space", `\s+`},
			{"ANY", `.`},
		},
	))
	parser = participle.MustBuild[File](
		participle.Lexer(lex),
		participle.Elide("Comment", "String", "MultilineString"),
		participle.UseLookahead(64),
	)
)

func Parse(filePath string) (*language.FileInfo, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	statements, err := parser.ParseBytes(filePath, content)
	if err != nil {
		return nil, err
	}
	return &language.FileInfo{
		Content: statements,
		Loc:     bytes.Count(content, []byte("\n")),
		Size:    len(content),
		AbsPath: filePath,
	}, nil
}
