package node

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNode_Level(t *testing.T) {
	tests := []struct {
		Name           string
		NumNodes       int
		Children       map[int][]int
		ExpectedLevels []int
	}{
		{
			Name:     "Simple",
			NumNodes: 4,
			Children: map[int][]int{
				0: {1, 2},
				1: {3},
				2: {3},
			},
			ExpectedLevels: []int{0, 1, 1, 2},
		},
		{
			Name:     "Cycle",
			NumNodes: 5,
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
			Name:     "Cycle 2",
			NumNodes: 3,
			Children: map[int][]int{
				0: {1, 2},
				1: {2, 0},
				2: {0, 1},
			},
			ExpectedLevels: []int{0, 2, 1},
		},
		{
			Name:     "Cycle 3",
			NumNodes: 3,
			Children: map[int][]int{
				0: {1},
				1: {2},
				2: {1},
			},
			ExpectedLevels: []int{0, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			nodes := make([]*Node[int], tt.NumNodes)
			for i := 0; i < tt.NumNodes; i++ {
				nodes[i] = MakeNode(strconv.Itoa(i), 0)
			}

			for n, children := range tt.Children {
				for _, child := range children {
					nodes[n].AddChild(nodes[child])
				}
			}
			ctx := context.Background()
			var lvls []int
			for i := 0; i < tt.NumNodes; i++ {
				var lvl int
				ctx, lvl = nodes[i].Level(ctx, nodes[0].Id)
				lvls = append(lvls, lvl)
			}
			a.Equal(tt.ExpectedLevels, lvls)
		})
	}
}
