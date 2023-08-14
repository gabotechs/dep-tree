//nolint:govet
package python_grammar

type VariableUnpack struct {
	Indented bool     `@Space?`
	Names    []string `( (@Ident (Space? "," Space? @Ident)+) | ( "(" (Space|NewLine)* @Ident ( (Space|NewLine)* "," (Space|NewLine)* @Ident )+ (Space|NewLine)* ")" ) ) Space? "="`
}

type VariableTyping struct {
	Indented bool   `@Space?`
	Name     string `@Ident Space? ":" Space? Ident`
}

type VariableAssign struct {
	Indented bool     `@Space?`
	Names    []string `(@Ident Space? "=" Space?)+`
}

type Function struct {
	Indented bool   `@Space?`
	Name     string `("async" Space)? "def" Space @Ident`
}

type Class struct {
	Indented bool   `@Space?`
	Name     string `"class" Space @Ident`
}
