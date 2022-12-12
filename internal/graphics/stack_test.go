package graphics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCellStack_Render_lines(t *testing.T) {
	a := require.New(t)
	cs := CellStack{}
	cs.add(LinesCell(Lines{
		t: true,
	}))
	cs.add(LinesCell(Lines{
		b: true,
	}))
	a.Equal("╷", cs.Render())
}
