package python_grammar

// This wraps a lexer so that it includes a SOF token with the start of the file

import (
	"io"
	"maps"

	"github.com/alecthomas/participle/v2/lexer"
)

func LexerWithSOF(parent lexer.Definition) lexer.Definition {
	out := maps.Clone(parent.Symbols())
	if _, ok := out["SOF"]; ok {
		panic("Wrapped lexer should not have the SOF token")
	}
	next := lexer.TokenType(-1)
	for _, t := range out {
		if t <= next {
			next = t - 1
		}
	}
	out["SOF"] = next - 1
	return &startLexerDef{parent: parent, symbols: out}
}

type startLexerDef struct {
	parent  lexer.Definition
	symbols map[string]lexer.TokenType
}

func (i *startLexerDef) Symbols() map[string]lexer.TokenType {
	return i.symbols
}

func (i *startLexerDef) Lex(filename string, r io.Reader) (lexer.Lexer, error) {
	lex, err := i.parent.Lex(filename, r)
	if err != nil {
		return nil, err
	}
	return &startLexer{
		lexer:     lex,
		startType: i.symbols["SOF"],
		started:   false,
	}, nil
}

var _ lexer.Definition = (*startLexerDef)(nil)

type startLexer struct {
	startType lexer.TokenType
	buffered  []lexer.Token
	lexer     lexer.Lexer
	started   bool
}

func (i *startLexer) Next() (lexer.Token, error) {
	if !i.started {
		i.started = true
		return lexer.Token{Pos: lexer.Position{}, Type: i.startType, Value: "^"}, nil
	}

	return i.lexer.Next()
}

var _ lexer.Lexer = (*startLexer)(nil)
