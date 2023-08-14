//nolint:govet
package python_grammar

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"

	"os"
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
	Statements []*Statement `(@@ | ANY | ALL | Punct | Ident | Space | NewLine)*`
	Path       string
}

var (
	lex = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"ALL", `\*`},
			{"Punct", `[,{}()]`},
			{"Ident", `[_$a-zA-Z][_$a-zA-Z0-9]*`},
			// https://stackoverflow.com/questions/69184441/regular-expression-for-comments-in-python-re
			{"Comment", `#.*`},
			{"NewLine", `\n+`},
			{"Space", `\s+`},
			{"ANY", `.`},
		},
	)
	parser = participle.MustBuild[File](
		participle.Lexer(lex),
		participle.Elide("Comment"),
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
