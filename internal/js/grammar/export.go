//nolint:govet
package grammar

type ExportDeconstruction struct {
	Names []string `"{" ((Ident "as" @Ident) | @Ident) ("," ((Ident "as" @Ident) | @Ident))* "}"`
}

type DeclarationExport struct {
	Name string `"export" "async"? ("let"|"const"|"var"|"function"|"class") ALL? @Ident`
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
