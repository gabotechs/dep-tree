//nolint:govet
package grammar

import (
	"context"
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
}

var (
	lex = lexer.MustSimple(
		[]lexer.SimpleRule{
			{"ALL", `\*`},
			{"Punct", `[,{}()]`},
			{"Ident", `[_$a-zA-Z\\xA0-\\uFFFF][_$a-zA-Z0-9\\xA0-\\uFFFF]*`},
			{"String", `'(\\'|[^'])*'|"(\\"|[^"])*"`},
			{"Comment", `//.*|/\*.*?\*/`},
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

type CacheKey string

func Parse(ctx context.Context, filePath string) (context.Context, *File, error) {
	if cached, ok := ctx.Value(CacheKey(filePath)).(*File); ok {
		return ctx, cached, nil
	} else {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return ctx, nil, err
		}
		file, err := parser.ParseBytes(filePath, content)
		if err != nil {
			return ctx, nil, err
		}
		ctx = context.WithValue(ctx, CacheKey(filePath), file)
		return ctx, file, nil
	}
}
