package graphics

import "dep-tree/internal/utils"

type CellStack []*TaggedCell

func (cs *CellStack) add(cell Cell) {
	if cs == nil {
		*cs = make([]*TaggedCell, 0)
	}
	*cs = append(*cs, NewTaggedCell(cell))
}

func (cs *CellStack) Tags() map[string]string {
	tags := map[string]string{}
	if cs != nil {
		for _, cell := range *cs {
			utils.Merge(tags, cell.tags)
		}
	}
	return tags
}

func (cs *CellStack) Is(key string, value string) bool {
	if cs == nil {
		return false
	}
	for _, cell := range *cs {
		if cell.Is(key, value) {
			return true
		}
	}
	return false
}

func (cs *CellStack) PlaceChar(char rune) *TaggedCell {
	charTaggedCell := NewTaggedCell(CharCell(char))
	cs.add(charTaggedCell)
	return charTaggedCell
}

func (cs *CellStack) PlaceArrow(inverted bool) *TaggedCell {
	arrowTaggedCell := NewTaggedCell(ArrowCell(inverted))
	cs.add(arrowTaggedCell)
	return arrowTaggedCell
}

func (cs *CellStack) PlaceEmpty() *TaggedCell {
	cell := NewTaggedCell(EmptyCell(false))
	cs.add(cell)
	return cell
}

var arrowMap = map[bool]rune{
	true:  '◁',
	false: '▷',
}

func (cs *CellStack) Render() rune {
	lineStack := LineStack{}

	for _, taggedCell := range *cs {
		switch cell := taggedCell.Cell.(type) {
		case EmptyCell:
			// nothing.
		case CharCell:
			return rune(cell)
		case ArrowCell:
			return arrowMap[bool(cell)]
		case *LinesCell:
			lineStack.add(cell)
		}
	}
	return lineStack.Render()
}
