//nolint:govet
package rust_grammar

type Mod struct {
	Pub   bool   `@"pub"? ("(" (Ident | PathSep)* ")")? "mod"`
	Name  string `@Ident`
	Local bool   `@"{"?`
}
