package graph

import (
	"sort"
	"strconv"
	"testing"

	"github.com/gabotechs/dep-tree/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestGraph_RemoveCycles(t *testing.T) {
	var tests = []struct {
		Name           string
		Children       [][]int
		ExpectedCauses [][2]int
		ExpectedCycles [][]int
	}{
		{
			Name: "Simple",
			Children: [][]int{
				0: {1, 2},
				1: {3},
				2: {3},
				3: {},
			},
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
			ExpectedCauses: [][2]int{{4, 3}},
			ExpectedCycles: [][]int{{3, 4, 3}},
		},
		{
			Name: "Cycle 2",
			Children: [][]int{
				0: {1, 2},
				1: {2, 0},
				2: {0, 1},
			},
			ExpectedCauses: [][2]int{{2, 1}, {2, 0}, {1, 0}},
			ExpectedCycles: [][]int{{1, 2, 1}, {0, 1, 2, 0}, {0, 1, 0}},
		},
		{
			Name: "Cycle 3",
			Children: [][]int{
				0: {1},
				1: {2},
				2: {1},
			},
			ExpectedCauses: [][2]int{{2, 1}},
			ExpectedCycles: [][]int{{1, 2, 1}},
		},
		{
			Name: "Cycle 4",
			Children: [][]int{
				0: {1},
				1: {2},
				2: {0},
			},
			ExpectedCauses: [][2]int{{2, 0}},
			ExpectedCycles: [][]int{{0, 1, 2, 0}},
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
			ExpectedCauses: [][2]int{{4, 1}},
			ExpectedCycles: [][]int{{1, 2, 3, 4, 1}},
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

			ExpectedCauses: [][2]int{{4, 2}},
			ExpectedCycles: [][]int{{2, 3, 4, 2}},
		},
		{
			Name: "Cycle 7",
			Children: [][]int{
				0: {1},
				1: {0, 2},
				2: {3},
				3: {4},
				4: {0},
			},
			ExpectedCauses: [][2]int{{4, 0}, {1, 0}},
			ExpectedCycles: [][]int{{0, 1, 2, 3, 4, 0}, {0, 1, 0}},
		},
		{
			Name: "Cycle 8",
			Children: [][]int{
				0: {3, 1},
				1: {2},
				2: {3},
				3: {4},
				4: {0},
			},
			ExpectedCauses: [][2]int{{4, 0}},
			ExpectedCycles: [][]int{{0, 3, 4, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			g := MakeTestGraph(tt.Children)

			cycles := g.RemoveCycles(g.Get("0"))

			var actualCycles [][]int
			var actualCauses [][2]int
			for _, c := range cycles {
				var actualCycle []int
				fromc, _ := strconv.Atoi(c.Cause[0])
				toc, _ := strconv.Atoi(c.Cause[1])
				actualCause := [2]int{fromc, toc}
				for _, el := range c.Stack {
					v, _ := strconv.Atoi(el)
					actualCycle = append(actualCycle, v)
				}
				actualCycles = append(actualCycles, actualCycle)
				actualCauses = append(actualCauses, actualCause)
			}
			sort.Slice(actualCycles, func(i, j int) bool {
				return utils.ItoAArr(actualCycles[i]) > utils.ItoAArr(actualCycles[j])
			})

			sort.Slice(actualCauses, func(i, j int) bool {
				return utils.ItoAArr2(actualCauses[i]) > utils.ItoAArr2(actualCauses[j])
			})
			a.Equal(tt.ExpectedCycles, actualCycles)
			a.Equal(tt.ExpectedCauses, actualCauses)
		})
	}
}
