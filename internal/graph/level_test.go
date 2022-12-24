package graph

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNode_Level(t *testing.T) {
	tests := []struct {
		Name           string
		Children       map[int][]int
		ExpectedLevels []int
	}{
		{
			Name: "Simple",
			Children: map[int][]int{
				0: {1, 2},
				1: {3},
				2: {3},
			},
			ExpectedLevels: []int{0, 1, 1, 2},
		},
		{
			Name: "Cycle",
			Children: map[int][]int{
				0: {1, 2, 3},
				1: {2, 4},
				2: {3, 4},
				3: {4},
				4: {3},
			},
			ExpectedLevels: []int{0, 1, 2, 4, 3},
		},
		{
			Name: "Cycle 2",
			Children: map[int][]int{
				0: {1, 2},
				1: {2, 0},
				2: {0, 1},
			},
			ExpectedLevels: []int{0, 2, 1},
		},
		{
			Name: "Cycle 3",
			Children: map[int][]int{
				0: {1},
				1: {2},
				2: {1},
			},
			ExpectedLevels: []int{0, 1, 2},
		},
		{
			Name: "Avoid same level",
			Children: map[int][]int{
				0: {1, 2},
				1: {},
				2: {1},
			},
			ExpectedLevels: []int{0, 2, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			numNodes := len(tt.ExpectedLevels)

			fr := NewGraph[int]()
			for i := 0; i < numNodes; i++ {
				fr.AddNode(MakeNode(strconv.Itoa(i), 0))
			}

			for n, children := range tt.Children {
				for _, child := range children {
					err := fr.AddChild(strconv.Itoa(n), strconv.Itoa(child))
					a.NoError(err)
				}
			}
			ctx := context.Background()
			var lvls []int
			for i := 0; i < numNodes; i++ {
				var lvl int
				ctx, lvl = fr.Level(ctx, strconv.Itoa(i), "0")
				lvls = append(lvls, lvl)
			}
			a.Equal(tt.ExpectedLevels, lvls)
		})
	}
}
