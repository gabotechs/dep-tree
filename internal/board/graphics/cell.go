package graphics

import "dep-tree/internal/utils"

type Cell interface {
	IsCell() bool
}

type LinesCell struct {
	l bool
	t bool
	r bool
	b bool
}

func (l *LinesCell) IsCell() bool {
	return true
}

type CharCell rune

func (c CharCell) IsCell() bool {
	return true
}

type ArrowCell bool

func (c ArrowCell) IsCell() bool {
	return true
}

type EmptyCell bool

func (c EmptyCell) IsCell() bool {
	return true
}

type TaggedCell struct {
	Cell
	tags map[string]string
}

func NewTaggedCell(cell Cell) *TaggedCell {
	switch t := cell.(type) {
	case *TaggedCell:
		return t
	default:
		return &TaggedCell{
			Cell: cell,
		}
	}
}

func (tc *TaggedCell) WithTags(tags map[string]string) *TaggedCell {
	if tc.tags == nil {
		tc.tags = tags
	} else {
		utils.Merge(tc.tags, tags)
	}
	return tc
}

func (tc *TaggedCell) WithTag(key string, value string) *TaggedCell {
	if tc.tags == nil {
		tc.tags = map[string]string{key: value}
	} else {
		tc.tags[key] = value
	}
	return tc
}

func (tc *TaggedCell) Is(key string, value string) bool {
	if tc.tags == nil {
		return false
	} else if v, ok := tc.tags[key]; ok {
		return value == v
	} else {
		return false
	}
}
