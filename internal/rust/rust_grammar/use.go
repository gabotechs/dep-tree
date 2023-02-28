//nolint:govet
package rust_grammar

type Name struct {
	Original string `@Ident`
	Alias    string `("as" @Ident)?`
}

type Use struct {
	Pub        bool     `@"pub"? ("(" (Ident | PathSep)* ")")?`
	PathSlices []string `"use" (@Ident PathSep)+`
	All        bool     `( @ALL |`
	Names      []Name   `         @@ | ( "{" @@ ("," @@)* "}" ) ) ";"`
}
