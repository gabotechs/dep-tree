//nolint:govet
package rust_grammar

type Pub struct {
	Name string `"pub"  ("(" (Ident | PathSep)* ")")? "unsafe"? ("fn" | "struct" | "trait" | "enum") @Ident`
}
