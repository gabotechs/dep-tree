package graphics

import "dep-tree/internal/utils"

type CellStack struct {
	cells []*Cell
}

func (c *CellStack) add(cell *Cell) {
	if c.cells == nil {
		c.cells = make([]*Cell, 0)
	}
	c.cells = append(c.cells, cell)
}

func (c *CellStack) Tags() map[string]string {
	tags := map[string]string{}
	if c.cells != nil {
		for _, cell := range c.cells {
			utils.Merge(tags, cell.tags)
		}
	}
	return tags
}

func (c *CellStack) Is(key string, value string) bool {
	if c.cells == nil {
		return false
	}
	for _, cell := range c.cells {
		if cell.Is(key, value) {
			return true
		}
	}
	return false
}

func (c *CellStack) PlaceChar(
	char rune,
) *Cell {
	cell := &Cell{
		t:    charCell,
		char: char,
	}
	c.add(cell)
	return cell
}

func (c *CellStack) PlaceArrow(
	inverted bool,
) *Cell {
	cell := &Cell{
		t:             arrowCell,
		arrowInverted: inverted,
	}
	c.add(cell)
	return cell
}

func (c *CellStack) PlaceEmpty() *Cell {
	cell := &Cell{
		t: emptyCell,
	}
	c.add(cell)
	return cell
}

var lineCharMap = map[int]rune{
	0b_0000: ' ',
	0b_0001: '╷',
	0b_0010: '╶',
	0b_0011: '┌',
	0b_0100: '╵',
	0b_0101: '│',
	0b_0110: '└',
	0b_0111: '├',
	0b_1000: '╴',
	0b_1001: '┐',
	0b_1010: '─',
	0b_1011: '┬',
	0b_1100: '┘',
	0b_1101: '┤',
	0b_1110: '┴',
	0b_1111: '┼',
}

var arrowMap = map[bool]rune{
	true:  '◁',
	false: '▷',
}

func hashLines(lines *Lines) int {
	result := 0b_0000
	if lines == nil {
		return result
	}
	if lines.l {
		result += 0b_1000
	}
	if lines.t {
		result += 0b_0100
	}
	if lines.r {
		result += 0b_0010
	}
	if lines.b {
		result += 0b_0001
	}
	return result
}

func (c *CellStack) Render() rune {
	var lines *Lines

	linesCrosses := false
	for _, cell := range c.cells {
		if cell.lines.cross {
			linesCrosses = true
		}
	}

	for _, cell := range c.cells {
		switch cell.t {
		case charCell:
			return cell.char
		case arrowCell:
			return arrowMap[cell.arrowInverted]
		case linesCell:
			if lines == nil || !linesCrosses {
				lines = &Lines{}
			}
			if cell.lines.l {
				lines.l = true
			}
			if cell.lines.t {
				lines.t = true
			}
			if cell.lines.r {
				lines.r = true
			}
			if cell.lines.b {
				lines.b = true
			}
		}
	}
	return lineCharMap[hashLines(lines)]
}
