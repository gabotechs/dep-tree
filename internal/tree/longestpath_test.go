package tree

import (
	"testing"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/stretchr/testify/require"
)

func Test_longestPath(t *testing.T) {
	var tests = []struct {
		Name           string
		Children       [][]int
		ExpectedLevels []int
		NoTrimCycles   bool
		ExpectedError  string
	}{
		{
			Name: "Simple",
			Children: [][]int{
				0: {1, 2},
				1: {3},
				2: {3},
				3: {},
			},
			ExpectedLevels: []int{0, 1, 1, 2},
		},
		{
			Name: "Cycle",
			Children: [][]int{
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
			Children: [][]int{
				0: {1, 2},
				1: {2, 0},
				2: {0, 1},
			},
			ExpectedLevels: []int{0, 1, 2},
		},
		{
			Name: "Cycle 3",
			Children: [][]int{
				0: {1},
				1: {2},
				2: {1},
			},
			ExpectedLevels: []int{0, 1, 2},
		},
		{
			Name: "Cycle 4",
			Children: [][]int{
				0: {1},
				1: {2},
				2: {0},
			},
			ExpectedLevels: []int{0, 1, 2},
		},
		{
			Name: "Cycle 5",
			Children: [][]int{
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
			Children: [][]int{
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
			Children: [][]int{
				0: {1, 2},
				1: {},
				2: {1},
			},
			ExpectedLevels: []int{0, 1, 2},
		},
		{
			Name: "Cycle (without trimming cycles first)",
			Children: [][]int{
				0: {1},
				1: {2},
				2: {1},
			},
			NoTrimCycles:  true,
			ExpectedError: "cannot calculate longest path between nodes because there is at least one cycle in the graph: cycle detected:\n1\n2\n1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			testParser := dep_tree.TestParser{
				Spec: tt.Children,
			}

			dt := dep_tree.NewDepTree[[]int](&testParser, []string{"0"})
			err := dt.LoadGraph()
			a.NoError(err)

			if !tt.NoTrimCycles {
				dt.LoadCycles()
			}

			numNodes := len(tt.Children)
			if tt.ExpectedError != "" {
				_, err = NewTree(dt)
				a.EqualError(err, tt.ExpectedError)
			} else {
				tree, err := NewTree(dt)
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
