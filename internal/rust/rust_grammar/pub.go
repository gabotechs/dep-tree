//nolint:govet
package rust_grammar

type Pub struct {
	Name Ident `"pub"  ("(" (Ident | PathSep)* ")")? "unsafe"? "async"? ("fn" | "struct" | "trait" | "enum" | "type" | "static" | "const") @Ident`
}
