package graphics

import "dep-tree/internal/utils"

const (
	charCell = iota
	linesCell
	arrowCell
	emptyCell
)

type Lines struct {
	l     bool
	t     bool
	r     bool
	b     bool
	cross bool
}

type Cell struct {
	t             int
	char          rune
	lines         Lines
	arrowInverted bool

	tags map[string]string
}

func (c *Cell) Tags(tags map[string]string) *Cell {
	if c.tags == nil {
		c.tags = tags
	} else {
		utils.Merge(c.tags, tags)
	}
	return c
}

func (c *Cell) Tag(key string, value string) *Cell {
	if c.tags == nil {
		c.tags = map[string]string{key: value}
	} else {
		c.tags[key] = value
	}
	return c
}

func (c *Cell) Is(key string, value string) bool {
	if c.tags == nil {
		return false
	} else if v, ok := c.tags[key]; ok {
		return value == v
	} else {
		return false
	}
}

func LinesCell(lines Lines) *Cell {
	return &Cell{
		t:     linesCell,
		lines: lines,
	}
}

func (c *Cell) mergeLines(lines Lines) {
	if lines.l {
		c.lines.l = true
	}
	if lines.t {
		c.lines.t = true
	}
	if lines.r {
		c.lines.r = true
	}
	if lines.b {
		c.lines.b = true
	}
	// NOTE: there is probably a better way of handling crossed lines.
	if (c.lines.t || c.lines.b) && (c.lines.r || c.lines.l) {
		c.lines.cross = true
	}
}
