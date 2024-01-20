//nolint:govet
package rust_grammar

import (
	"strings"
)

type Ident string

func (i *Ident) Capture(values []string) error {
	if strings.HasPrefix(values[0], "r#") {
		*i = Ident(values[0][2:])
	} else {
		*i = Ident(values[0])
	}
	return nil
}

type Mod struct {
	Pub   bool  `@"pub"? ("(" (Ident | PathSep)* ")")? "mod"`
	Name  Ident `@Ident`
	Local bool  `@"{"?`
}
