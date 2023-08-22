package graphics

import "github.com/gabotechs/dep-tree/internal/utils"

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

func (cs *CellStack) Tag(key string) string {
	if cs == nil {
		return ""
	}
	for _, cell := range *cs {
		tag := cell.Tag(key)
		if tag != "" {
			return tag
		}
	}
	return ""
}

func (cs *CellStack) Match(tags map[string]string) bool {
	if tags == nil {
		return false
	}
	for k, v := range tags {
		if cs.Is(k, v) {
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

func (cs *CellStack) Render(
	priorityTags map[string]string,
) rune {
	lineStack := LineStack{}
	priorityCellStack := *cs

	if cs.Match(priorityTags) {
		priorityCellStack = CellStack{}
		for _, cell := range *cs {
			if cell.Match(priorityTags) {
				priorityCellStack.add(cell)
			}
		}
	}

	for _, taggedCell := range priorityCellStack {
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
