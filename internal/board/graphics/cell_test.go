package graphics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCell_mergeLines(t *testing.T) {
	a := require.New(t)
	cell := LinesCell(Lines{
		t: true,
	})
	cell.mergeLines(Lines{
		b: true,
	})
	a.Equal(cell.lines.l, false)
	a.Equal(cell.lines.t, true)
	a.Equal(cell.lines.r, false)
	a.Equal(cell.lines.b, true)
	a.Equal(cell.lines.cross, false)

	cell.mergeLines(Lines{
		l: true,
	})

	a.Equal(cell.lines.l, true)
	a.Equal(cell.lines.t, true)
	a.Equal(cell.lines.r, false)
	a.Equal(cell.lines.b, true)
	a.Equal(cell.lines.cross, true)
}

func TestCell_Tag(t *testing.T) {
	a := require.New(t)
	c := Cell{}

	a.Equal(false, c.Is("key", "foo"))
	c.Tag("key", "bar")
	a.Equal(false, c.Is("key", "foo"))
	a.Equal(true, c.Is("key", "bar"))
	c.Tag("key", "foo")
	c.Tag("otherKey", "bar")
	a.Equal(true, c.Is("key", "foo"))
	a.Equal(true, c.Is("otherKey", "bar"))
}
