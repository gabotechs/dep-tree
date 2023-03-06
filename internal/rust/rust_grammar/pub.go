//nolint:govet
package rust_grammar

type Pub struct {
	Name string `"pub"  ("(" (Ident | PathSep)* ")")? "unsafe"? "async"? ("fn" | "struct" | "trait" | "enum" | "type" | "static") @Ident`
}
