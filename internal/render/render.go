package render

import "golang.org/x/text/unicode/norm"

var lineCharMap = map[int]rune{
	0b_0000: ' ',
	0b_0001: '╷',
	0b_0010: '╶',
	0b_0011: '┌',
	0b_0100: '╵',
	0b_0101: '│',
	0b_0110: '└',
	0b_0111: '├',
	0b_1000: ' ', // lost this one
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
			if lines == nil {
				lines = &Lines{}
			} else if !linesCrosses {
				continue
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
