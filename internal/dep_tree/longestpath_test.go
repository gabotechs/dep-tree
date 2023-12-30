package dep_tree

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	testgraph "github.com/gabotechs/dep-tree/internal/graph"
)

func Test_longestPath(t *testing.T) {
	var tests = []struct {
		Name           string
		Children       map[int][]int
		ExpectedLevels []int
		NoTrimCycles   bool
		ExpectedError  string
	}{
		{
			Name: "Simple",
			Children: map[int][]int{
				0: {1, 2},
				1: {3},
				2: {3},
				3: {},
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
			ExpectedLevels: []int{0, 1, 2, 3, 4},
		},
		{
			Name: "Cycle 2",
			Children: map[int][]int{
				0: {1, 2},
				1: {2, 0},
				2: {0, 1},
			},
			ExpectedLevels: []int{0, 1, 2},
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
			Name: "Cycle 4",
			Children: map[int][]int{
				0: {1},
				1: {2},
				2: {0},
			},
			ExpectedLevels: []int{0, 1, 2},
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
			ExpectedLevels: []int{0, 1, 2, 3, 4},
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
			ExpectedLevels: []int{0, 1, 2, 3, 4},
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
		{
			Name: "Cycle (without trimming cycles first)",
			Children: map[int][]int{
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

			g := testgraph.MakeTestGraph(tt.Children)

			var err error
			if !tt.NoTrimCycles {
				g.RemoveCycles(g.Get("0"))
				a.NoError(err)
			}

			numNodes := len(tt.Children)
			dt := DepTree[int]{longestPathCache: map[string]int{}}
			if tt.ExpectedError != "" {
				for i := 0; i < numNodes; i++ {
					_, err = dt.longestPath(g, "0", strconv.Itoa(i), nil)
					if err != nil {
						break
					}
				}
				a.EqualError(err, tt.ExpectedError)
			} else {
				var lvls []int
				for i := 0; i < numNodes; i++ {
					var lvl int
					lvl, err = dt.longestPath(g, "0", strconv.Itoa(i), nil)
					a.NoError(err)
					lvls = append(lvls, lvl)
				}
				a.Equal(tt.ExpectedLevels, lvls)
			}
		})
	}
}
