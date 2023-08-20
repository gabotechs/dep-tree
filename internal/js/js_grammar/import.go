//nolint:govet
package js_grammar

type ImportDeconstruction struct {
	Names []string `"{" "type"? @Ident ("as" Ident)? ("," ("type"? @Ident ("as" Ident)?)?)* "}"`
}

type AllImport struct {
	Alias string `ALL ("as" @Ident)?`
}

type SelectionImport struct {
	AllImport      *AllImport            `(@@?`
	Deconstruction *ImportDeconstruction ` @@?)!`
}

type Imported struct {
	Default         bool             `(@Ident? ","?`
	SelectionImport *SelectionImport ` @@?)!`
}

type StaticImport struct {
	Imported *Imported `"import" "type"? (@@ "from")?`
	Path     string    `@String`
}

type DynamicImport struct {
	Path string `"import" "(" @String ")"`
}
