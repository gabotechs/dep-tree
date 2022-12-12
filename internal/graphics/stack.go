package graphics

import "golang.org/x/text/unicode/norm"

type CellStack struct {
	tags  map[string]string
	cells []*Cell
}

func (c *CellStack) add(cell *Cell) {
	if c.cells == nil {
		c.cells = make([]*Cell, 0)
	}
	c.cells = append(c.cells, cell)
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

func (c *CellStack) Render() string {
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
			result := ""
			for _, r := range cell.char.runes {
				result += string(r)
			}
			return norm.NFC.String(result)
		case arrowCell:
			return string(arrowMap[cell.arrowInverted])
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
	return string(lineCharMap[hashLines(lines)])
}
