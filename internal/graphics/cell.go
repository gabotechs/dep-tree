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

func LinesCell(lines Lines) *Cell {
	return &Cell{
		t:     linesCell,
		lines: lines,
	}
}

func (c *Cell) mergeLines(lines Lines) {
	if lines.l {
		c.lines.l = true
	}
	if lines.t {
		c.lines.t = true
	}
	if lines.r {
		c.lines.r = true
	}
	if lines.b {
		c.lines.b = true
	}
	// NOTE: there is probably a better way of handling crossed lines.
	if (c.lines.t || c.lines.b) && (c.lines.r || c.lines.l) {
		c.lines.cross = true
	}
}
