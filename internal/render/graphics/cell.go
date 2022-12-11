package graphics

const (
	charCell = iota
	linesCell
	arrowCell
)

type Lines struct {
	l     bool
	t     bool
	r     bool
	b     bool
	cross bool
}

type Char struct {
	runes []rune
}

type Cell struct {
	t             int
	char          Char
	lines         Lines
	arrowInverted bool
}
