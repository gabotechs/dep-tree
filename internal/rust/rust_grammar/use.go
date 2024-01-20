//nolint:govet
package rust_grammar

type Name struct {
	Original Ident `@Ident`
	Alias    Ident `("as" @Ident)?`
}

type UsePath struct {
	PathSlices []Ident    `(@Ident PathSep)*`
	All        bool       `(@ALL |`
	Name       *Name      `      @@ |`
	UsePaths   []*UsePath `         "{" @@ ("," @@)* "}" )`
}

type Use struct {
	Pub     bool     `@"pub"? ("(" (Ident | PathSep)* ")")?`
	UsePath *UsePath `"use" @@ ";"`
}

type FlattenUse struct {
	Pub        bool
	PathSlices []string
	All        bool
	Name       Name
}

func flatten(node *UsePath, pathSlices []string, isPub bool) []FlattenUse {
	currentSlices := pathSlices
	for _, pathSlice := range node.PathSlices {
		currentSlices = append(currentSlices, string(pathSlice))
	}
	switch {
	case node.All:
		return []FlattenUse{
			{
				Pub:        isPub,
				PathSlices: currentSlices,
				All:        true,
			},
		}
	case node.Name != nil:
		return []FlattenUse{
			{
				Pub:        isPub,
				PathSlices: currentSlices,
				Name:       *node.Name,
			},
		}
	default:
		uses := make([]FlattenUse, 0)
		for _, use := range node.UsePaths {
			uses = append(uses, flatten(use, currentSlices, isPub)...)
		}
		return uses
	}
}

func (u *Use) Flatten() []FlattenUse {
	return flatten(u.UsePath, []string{}, u.Pub)
}
