package graphics

type CellStack struct {
	tags  map[string]string
	cells []*Cell
}

func (c *CellStack) add(cell *Cell) {
	c.cells = append([]*Cell{cell}, c.cells...)
}

func (c *CellStack) Tag(key string, value string) {
	if c.tags == nil {
		c.tags = map[string]string{key: value}
	} else {
		c.tags[key] = value
	}
}

func (c *CellStack) Is(key string, value string) bool {
	if c.tags == nil {
		return false
	} else if v, ok := c.tags[key]; ok {
		return value == v
	} else {
		return false
	}
}

func (c *CellStack) PlaceChar(chars ...rune) bool {
	c.add(&Cell{
		t: charCell,
		char: Char{
			runes: chars,
		},
	})
	return true
}

func (c *CellStack) PlaceArrow(inverted bool) bool {
	c.add(&Cell{
		t:             arrowCell,
		arrowInverted: inverted,
	})
	return true
}
