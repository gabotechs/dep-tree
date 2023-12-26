//nolint:govet
package js_grammar

type AliasedName struct {
	Original string `@Ident`
	Alias    string `("as" @Ident)?`
}

type ExportDeconstruction struct {
	Names []AliasedName `"{" @@ ("," @@)* ","? "}"`
}

type DeclarationExport struct {
	Name string `"export" "async"? ("let"|"const"|"var"|"function"|"class"|"type"|"interface"|"enum") ALL? @Ident`
}

type ListExport struct {
	ExportDeconstruction *ExportDeconstruction `"export" @@`
}

type DefaultExport struct {
	Default bool `"export" @"default"`
}

type ProxyExport struct {
	ExportDeconstruction *ExportDeconstruction `"export" (@@`
	ExportAll            bool                  `             | (@ALL`
	ExportAllAlias       string                `                     ("as" @Ident)?)) `
	From                 string                `"from" @String`
}
