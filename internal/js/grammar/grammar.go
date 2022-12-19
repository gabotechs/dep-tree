//nolint:govet
package grammar

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"

	"dep-tree/internal/utils"
)

type Statement struct {
	// imports
	DynamicImport *DynamicImport `@@`
	StaticImport  *StaticImport  `| @@`
	// exports
	DeclarationExport *DeclarationExport `| @@`
	ProxyExport       *ProxyExport       `| @@`
	DefaultExport     *DefaultExport     `| @@`
	ListExport        *ListExport        `| @@`
}

type File struct {
	Statements []*Statement `(@@ | ANY | FALSE_IMPORT_1 | FALSE_IMPORT_2 | FALSE_EXPORT_1 | FALSE_EXPORT_2)*`
}

const identRe = `[_$a-zA-Z\\xA0-\\uFFFF][_$a-zA-Z0-9\\xA0-\\uFFFF]*`
const stringRe = `'[^']*'|"[^"]*"`
const commentRe = `//.*|/\*.*?\*/`

const punctuationRe = `[,{}()]`

var (
	lex = lexer.MustStateful(utils.Merge(
		lexer.Rules{
			"Root": {
				{"FALSE_IMPORT_1", `import` + identRe, nil},
				{"FALSE_IMPORT_2", identRe + `import`, nil},
				{"IMPORT", `import`, lexer.Push("Import")},
				{"FALSE_EXPORT_1", `export` + identRe, nil},
				{"FALSE_EXPORT_2", identRe + `export`, nil},

				{"DECLARATION_EXPORT", `export\s+(const|let|var|function|class)`, lexer.Push("DeclarationExport")},
				{"PROXY_EXPORT", `export\s*\{[\s,_$a-zA-Z0-9\\xA0-\\uFFFF]*}\s*from`, lexer.Push("ProxyExport")},
				{"LIST_EXPORT", `export\s*{`, lexer.Push("ListExport")},
				{"DEFAULT_EXPORT", `export\s+default`, lexer.Push("DefaultExport")},

				{"Whitespace", `\s+`, nil},
				{"ANY", `.`, nil},
			},
		},
		importLexer,
		exportLexer,
	))
	parser = participle.MustBuild[File](
		participle.Lexer(lex),
		participle.Elide("Whitespace", "Comment"),
		participle.Unquote("String"),
		participle.UseLookahead(1024),
	)
)

func Parse(content []byte) (*File, error) {
	return parser.ParseBytes("", content)
}
