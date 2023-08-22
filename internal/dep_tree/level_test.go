package dep_tree

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/graph"
)

func TestNode_Level(t *testing.T) {
	tests := []struct {
		Name           string
		Children       map[int][]int
		ExpectedLevels []int
		ExpectedCycles [][]string
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
			ExpectedCycles: [][]string{{"3", "4", "3"}},
		},
		{
			Name: "Cycle 2",
			Children: map[int][]int{
				0: {1, 2},
				1: {2, 0},
				2: {0, 1},
			},
			ExpectedLevels: []int{0, 2, 1},
			ExpectedCycles: [][]string{{"1", "2", "1"}},
		},
		{
			Name: "Cycle 3",
			Children: map[int][]int{
				0: {1},
				1: {2},
				2: {1},
			},
			ExpectedLevels: []int{0, 1, 2},
			ExpectedCycles: [][]string{{"1", "2", "1"}},
		},
		{
			Name: "Cycle 4",
			Children: map[int][]int{
				0: {1},
				1: {2},
				2: {0},
			},
			ExpectedLevels: []int{0, 1, 2},
			// TODO: there is a bug with cyclical deps belonging to node 0.
			ExpectedCycles: nil,
		},
		{
			Name: "Cycle 5",
			Children: map[int][]int{
				0: {1},
				1: {2},
				2: {3},
				3: {4},
				4: {1},
			},
			// TODO: weird level calculation, but strictly correct.
			ExpectedLevels: []int{0, 3, 4, 1, 2},
			ExpectedCycles: [][]string{{"1", "2", "3", "4", "1"}},
		},
		{
			Name: "Cycle 6",
			Children: map[int][]int{
				0: {1},
				1: {2},
				2: {3},
				3: {4},
				4: {2},
			},
			ExpectedLevels: []int{0, 1, 2, 3, 1},
			ExpectedCycles: [][]string{{"2", "3", "4", "2"}},
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

			g := graph.NewGraph[int]()
			for i := 0; i < numNodes; i++ {
				g.AddNode(graph.MakeNode(strconv.Itoa(i), 0))
			}

			for n, children := range tt.Children {
				for _, child := range children {
					err := g.AddChild(strconv.Itoa(n), strconv.Itoa(child))
					a.NoError(err)
				}
			}
			levelCalculator := NewLevelCalculator(g, "0")

			ctx := context.Background()
			var lvls []int
			for i := 0; i < numNodes; i++ {
				var lvl int
				ctx, lvl = levelCalculator.level(ctx, strconv.Itoa(i))
				lvls = append(lvls, lvl)
			}
			a.Equal(tt.ExpectedLevels, lvls)
			var cycles [][]string
			for _, cycleKey := range levelCalculator.Cycles.Keys() {
				cycle, _ := levelCalculator.Cycles.Get(cycleKey)
				cycles = append(cycles, cycle.Stack)
			}
			a.Equal(tt.ExpectedCycles, cycles)
		})
	}
}
