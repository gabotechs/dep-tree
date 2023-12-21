package utils

import (
	"strconv"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// UnquoteSafe applies strconv.Unquote() to tokens of the given types.
//
// Tokens of type "String" will be unquoted if no other types are provided.
func UnquoteSafe(types ...string) participle.Option {
	if len(types) == 0 {
		types = []string{"String"}
	}
	return participle.Map(func(t lexer.Token) (lexer.Token, error) {
		value, err := unquote(t.Value)
		if err != nil {
			return t, nil
		}
		t.Value = value
		return t, nil
	}, types...)
}

func unquote(s string) (string, error) {
	quote := s[0]
	s = s[1 : len(s)-1]
	out := ""
	for s != "" {
		value, _, tail, err := strconv.UnquoteChar(s, quote)
		if err != nil {
			return "", err
		}
		s = tail
		out += string(value)
	}
	return out, nil
}
