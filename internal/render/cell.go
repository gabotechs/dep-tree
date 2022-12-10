package render

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

type CellStack struct {
	cells []Cell
}

func (c *CellStack) add(cell Cell) {
	c.cells = append([]Cell{cell}, c.cells...)
}

func (c *CellStack) PlaceChar(chars ...rune) {
	c.add(Cell{
		t: charCell,
		char: Char{
			runes: chars,
		},
	})
}

func (c *CellStack) DrawVerticalLine() {
	c.add(Cell{
		t: linesCell,
		lines: Lines{
			t: true,
			b: true,
		},
	})
}

func (c *CellStack) DrawHorizontalLine() {
	c.add(Cell{
		t: linesCell,
		lines: Lines{
			l: true,
			r: true,
		},
	})
}

// DrawJoint
//
//	true,  true  -> "┌"
//	true,  false -> "└"
//	false, true  -> "┐"
//	false, false -> "┘"
func (c *CellStack) DrawJoint(x bool, y bool) {
	c.add(Cell{
		t: linesCell,
		lines: Lines{
			l:     !x,
			t:     !y,
			r:     x,
			b:     y,
			cross: true,
		},
	})
}

func (c *CellStack) PlaceArrow(inverted bool) {
	c.add(Cell{
		t:             arrowCell,
		arrowInverted: inverted,
	})
}
