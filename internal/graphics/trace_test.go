package graphics

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"dep-tree/internal/vector"
)

const (
	up    = "up"
	down  = "down"
	left  = "left"
	right = "right"
)

func TestLineTracer(t *testing.T) {
	tests := []struct {
		Name     string
		Expected string
		StartX   int
		StartY   int
	}{
		{
			Name: "down,down,right,down",
			Expected: `
╷    
│    
└┐   
 ╵   
     
`,
		},
		{
			Name: "right,right,right,right,down,down,down,down,left,left,left",
			Expected: `
╶───┐
    │
    │
    │
 ╶──┘
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			dirs := strings.Split(tt.Name, ",")

			matrix := NewMatrix(5, 5)

			position := vector.Vec(tt.StartX, tt.StartY)
			tracer := NewLineTracer(position)
			for _, dir := range dirs {
				switch dir {
				case up:
					tracer.MoveVertical(true)
				case down:
					tracer.MoveVertical(false)
				case left:
					tracer.MoveHorizontal(true)
				case right:
					tracer.MoveHorizontal(false)
				}
			}

			err := tracer.Dump(matrix)
			a.NoError(err)
			result := matrix.Render()
			a.Equal(strings.TrimLeft(tt.Expected, " \n"), result)
		})
	}
}
