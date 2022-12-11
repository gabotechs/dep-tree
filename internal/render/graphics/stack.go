package graphics

type CellStack struct {
	IsSolid bool
	cells   []Cell
}

func NewCellStack() CellStack {
	return CellStack{
		cells: []Cell{},
	}
}

func (c *CellStack) add(cell Cell) {
	c.cells = append([]Cell{cell}, c.cells...)
}

func (c *CellStack) solid() {
	c.IsSolid = true
}

func (c *CellStack) PlaceChar(chars ...rune) {
	c.solid()
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
