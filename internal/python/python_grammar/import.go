//nolint:govet
package python_grammar

type ImportedName struct {
	Name  string `@Ident`
	Alias string `(Space "as" Space @Ident)?`
}

type FromImport struct {
	Indented bool           `@Space?`
	Relative []bool         `"from" (Space | @".") @"."*`
	Path     []string       `@Ident? (Space? "." Space? @Ident)* Space`
	All      bool           `"import" Space ( @ALL |`
	Names    []ImportedName `( (@@ (Space? "," Space? @@)*) | ( "(" (Space|NewLine)* @@ ( (Space|NewLine)* "," (Space|NewLine)* @@ )* (Space|NewLine)* ","? (Space|NewLine)* ")" ) ) )`
}

type Import struct {
	Indented bool     `@Space?`
	Path     []string `"import" Space @Ident (Space? "." Space? @Ident)*`
	Alias    string   `(Space "as" Space @Ident)?`
}
