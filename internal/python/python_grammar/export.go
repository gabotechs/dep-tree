//nolint:govet
package python_grammar

type VariableUnpack struct {
	Names []string `(NewLine | SOF) ( (@Ident (Space? "," Space? @Ident)+) | ( "(" (Space|NewLine)* @Ident ( (Space|NewLine)* "," (Space|NewLine)* @Ident )+ (Space|NewLine)* ")" ) ) Space? "="`
}

type VariableTyping struct {
	Name string `(NewLine | SOF) @Ident Space? ":" Space? Ident`
}

type VariableAssign struct {
	Names []string `(NewLine | SOF) (@Ident Space? "=" Space?)+`
}

type Function struct {
	Name string `(NewLine | SOF) ("async" Space)? "def" Space @Ident`
}

type Class struct {
	Name string `(NewLine | SOF) "class" Space @Ident`
}
