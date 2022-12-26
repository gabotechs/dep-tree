package graphics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLineStack_Render(t *testing.T) {
	a := require.New(t)
	lineStack := LineStack{}
	lineStack.add(&LinesCell{t: true, b: true})
	a.Equal('│', lineStack.Render())
	lineStack.add(&LinesCell{b: true})
	a.Equal('│', lineStack.Render())
	lineStack.add(&LinesCell{l: true})
	a.Equal('┤', lineStack.Render())
}

func TestLineStack_Render_Crossed(t *testing.T) {
	a := require.New(t)
	lineStack := LineStack{}
	lineStack.add(&LinesCell{t: true, b: true})
	lineStack.add(&LinesCell{l: true, r: true})
	a.Equal('─', lineStack.Render())
	lineStack.add(&LinesCell{t: true, b: true})
	a.Equal('│', lineStack.Render())
}

func TestLineStack_Render_Crossed_2(t *testing.T) {
	a := require.New(t)
	lineStack := LineStack{}
	lineStack.add(&LinesCell{t: true, r: true})
	lineStack.add(&LinesCell{l: true, r: true})
	a.Equal('┴', lineStack.Render())
	lineStack.add(&LinesCell{t: true, b: true})
	a.Equal('│', lineStack.Render())
}
