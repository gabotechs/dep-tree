package ecs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type CorrectComponent int

type AlternativeComponent struct{}

type NonPresentComponent struct{}

func TestSystem(t *testing.T) {
	tests := []struct {
		Name         string
		Systems      []interface{}
		ExpectedData int
	}{
		{
			Name: "Non function system",
			Systems: []interface{}{
				1,
				"string",
			},
		},
		{
			Name: "System referencing non present component",
			Systems: []interface{}{
				func(c *NonPresentComponent) {},
				func(c *CorrectComponent) { *c++ },
				func(c *CorrectComponent, c2 *AlternativeComponent) { *c++ },
				func(c *CorrectComponent, c2 *NonPresentComponent) { *c++ },
				func(c *CorrectComponent, c2 *AlternativeComponent, c3 *NonPresentComponent) { *c++ },
			},
			ExpectedData: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			component := CorrectComponent(0)

			w := NewWorld().
				WithEntity(NewEntity().
					With(&component).
					With(&AlternativeComponent{}),
				)
			for _, system := range tt.Systems {
				w = w.WithSystem(system)
			}

			err := w.Update()
			a.NoError(err)
			a.Equal(CorrectComponent(tt.ExpectedData), component)
		})
	}
}
