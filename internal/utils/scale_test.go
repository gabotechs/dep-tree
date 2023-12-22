package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScale(t *testing.T) {
	tests := []struct {
		Name     string
		V        [5]float64
		Expected float64
	}{
		{
			Name:     "1",
			V:        [5]float64{1, 0.5, 1.5, 0, 2},
			Expected: 1,
		},
		{
			Name:     "2",
			V:        [5]float64{-2, 0.5, 1.5, 0, 2},
			Expected: 0,
		},
		{
			Name:     "3",
			V:        [5]float64{10, 0.5, 1.5, 0, 2},
			Expected: 2,
		},
		{
			Name:     "4",
			V:        [5]float64{.75, 0.5, 1, 1, 2},
			Expected: 1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			a.Equal(tt.Expected, Scale(tt.V[0], tt.V[1], tt.V[2], tt.V[3], tt.V[4]))
		})
	}
}
