package tree

import (
	"testing"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/stretchr/testify/require"
)

func Test_longestPath(t *testing.T) {
	var tests = []struct {
		Name           string
		Spec           [][]int
		ExpectedLevels []int
		ExpectedError  string
	}{
		{
			Name: "Simple",
			Spec: [][]int{
				0: {1, 2},
				1: {3},
				2: {3},
				3: {},
			},
			ExpectedLevels: []int{0, 1, 1, 2},
		},
		{
			Name: "Cycle",
			Spec: [][]int{
				0: {1, 2, 3},
				1: {2, 4},
				2: {3, 4},
				3: {4},
				4: {3},
			},
			ExpectedLevels: []int{0, 1, 2, 3, 4},
		},
		{
			Name: "Cycle 2",
			Spec: [][]int{
				0: {1, 2},
				1: {2, 0},
				2: {0, 1},
			},
			ExpectedLevels: []int{0, 1, 2},
		},
		{
			Name: "Cycle 3",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {1},
			},
			ExpectedLevels: []int{0, 1, 2},
		},
		{
			Name: "Cycle 4",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {0},
			},
			ExpectedLevels: []int{0, 1, 2},
		},
		{
			Name: "Cycle 5",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {3},
				3: {4},
				4: {1},
			},
			ExpectedLevels: []int{0, 1, 2, 3, 4},
		},
		{
			Name: "Cycle 6",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {3},
				3: {4},
				4: {2},
			},
			ExpectedLevels: []int{0, 1, 2, 3, 4},
		},
		{
			Name: "Avoid same level",
			Spec: [][]int{
				0: {1, 2},
				1: {},
				2: {1},
			},
			ExpectedLevels: []int{0, 1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			tree, err := NewTree[[]int](
				[]string{"0"},
				&graph.TestParser{Spec: tt.Spec},
				func(node *graph.Node[[]int]) string { return node.Id },
				nil,
			)
			numNodes := len(tt.Spec)
			if tt.ExpectedError != "" {
				a.EqualError(err, tt.ExpectedError)
			} else {
				a.NoError(err)
				var lvls []int
				for i := 0; i < numNodes; i++ {
					lvls = append(lvls, tree.Nodes[i].Lvl)
				}
				a.Equal(tt.ExpectedLevels, lvls)
			}
		})
	}
}
