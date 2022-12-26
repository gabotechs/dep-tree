package graphics

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

func hashLines(lines *LinesCell) int {
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

type LineStack []*LinesCell

func (l *LineStack) add(cell *LinesCell) {
	*l = append(*l, cell)
}

func areCrossing(first *LinesCell, last *LinesCell) bool {
	if first == nil || last == nil {
		return false
	}
	if hashLines(first) == 0b_0101 && hashLines(last) == 0b_1010 {
		// lines have crossed.
		return true
	} else if hashLines(last) == 0b_0101 && hashLines(first) == 0b_1010 {
		// lines have crossed.
		return true
	}
	return false
}

func (l *LineStack) Render() rune {
	acc := LinesCell{}
	var prev *LinesCell
	for _, cell := range *l {
		// TODO: honestly, there should be a more elegant way of handling crossed lines.
		if areCrossing(prev, cell) {
			acc = LinesCell{}
		}
		if cell.l {
			acc.l = true
		}
		if cell.t {
			acc.t = true
		}
		if cell.r {
			acc.r = true
		}
		if cell.b {
			acc.b = true
		}
		prev = cell
	}
	return lineCharMap[hashLines(&acc)]
}
