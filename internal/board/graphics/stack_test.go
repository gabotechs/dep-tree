package graphics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCellStack_Render_lines(t *testing.T) {
	a := require.New(t)
	cs := CellStack{}
	cs.add(NewTaggedCell(&LinesCell{
		t: true,
	}))
	cs.add(NewTaggedCell(&LinesCell{
		b: true,
	}))
	a.Equal('│', cs.Render(nil))
}

func TestCellStack_Render_linesPriority(t *testing.T) {
	a := require.New(t)
	cs := CellStack{}
	cs.add(NewTaggedCell(&LinesCell{
		t: true,
		r: true,
	}))
	cs.add(NewTaggedCell(&LinesCell{
		l: true,
		r: true,
	}).WithTag("test", "true"))
	a.Equal('─', cs.Render(map[string]string{"test": "true"}))
}

func TestCellStack_Render_charHasPriority(t *testing.T) {
	a := require.New(t)
	cs := CellStack{}
	cs.add(NewTaggedCell(&LinesCell{
		t: true,
	}))
	cs.PlaceChar('a')
	a.Equal('a', cs.Render(nil))
}

func TestCellStack_Render_arrowHasPriority(t *testing.T) {
	a := require.New(t)
	cs := CellStack{}
	cs.PlaceArrow(false)
	cs.add(NewTaggedCell(&LinesCell{
		t: true,
	}))
	a.Equal('▷', cs.Render(nil))
}
