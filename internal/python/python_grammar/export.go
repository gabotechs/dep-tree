//nolint:govet
package python_grammar

type Variable struct {
	Indented bool   `@Space?`
	Name     string `@Ident Space? ("=" | (":" Space? Ident))`
}

type Function struct {
	Indented bool   `@Space?`
	Name     string `("async" Space)? "def" Space @Ident`
}

type Class struct {
	Indented bool   `@Space?`
	Name     string `"class" Space @Ident`
}
