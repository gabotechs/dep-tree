package graph

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadDeps_noOwnChild(t *testing.T) {
	a := require.New(t)
	testParser := TestParser{
		Spec: [][]int{{0}},
	}
	g := NewGraph[[]int]()
	err := g.Load([]string{"0"}, &testParser, nil)
	a.NoError(err)

	a.Equal(0, len(g.ToId("0")))
}

func TestLoadDeps_ErrorHandle(t *testing.T) {
	a := require.New(t)
	testParser := TestParser{
		Spec: [][]int{
			{1},
			{2},
			{-3},
		},
	}
	g := NewGraph[[]int]()
	err := g.Load([]string{"0"}, &testParser, nil)
	a.NoError(err)

	g.RemoveCycles([]*Node[[]int]{MakeNode("0", testParser.Spec[0])})

	node0 := g.Get("0")
	a.NotNil(node0)
	a.Equal(len(node0.Errors), 0)
	node1 := g.Get("1")
	a.NotNil(node1)
	a.Equal(len(node1.Errors), 0)
	node2 := g.Get("2")
	a.NotNil(node2)
	a.ErrorContains(node2.Errors[0], "no negative children")
}

func TestLoadDeps_Callbacks(t *testing.T) {
	a := require.New(t)

	testParser := TestParser{
		Spec: [][]int{
			0: {1},
			1: {2},
			2: {0, 3},
			3: {4},
			4: {3},
			5: {},
		},
	}

	testCallbacks := TestCallbacks[[]int]{}

	g := NewGraph[[]int]()
	err := g.Load([]string{"0"}, &testParser, &testCallbacks)
	a.NoError(err)

	a.Equal(1, testCallbacks.startLoad)
	a.Equal(1, testCallbacks.finishLoad)
	a.Equal(5, testCallbacks.nodeLoaded)
}

func TestLoadDeps_loadGraph(t *testing.T) {
	tests := []struct {
		Name        string
		Spec        [][]int
		Ids         []int
		Entrypoints []int
		NCycles     int
	}{
		{
			Name: "inner cycle",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {3},
				3: {4},
				4: {2},
			},
			Ids:         []int{0, 1, 2, 3, 4},
			Entrypoints: []int{0},
			NCycles:     1,
		},
		{
			Name: "cycle from 0",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {3},
				3: {4},
				4: {0},
			},
			Ids:         []int{0, 1, 2, 3, 4},
			Entrypoints: []int{},
			NCycles:     1,
		},
		{
			Name: "two cycles",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {3, 1},
				3: {4},
				4: {3},
			},
			Ids:         []int{0, 1, 2, 3, 4},
			Entrypoints: []int{0},
			NCycles:     2,
		},
		{
			Name: "three cycles, one from 0",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {3, 1},
				3: {4},
				4: {3, 0},
			},
			Ids:         []int{0, 1, 2, 3, 4},
			Entrypoints: []int{},
			NCycles:     3,
		},
		{
			Name: "two clusters",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {},
				3: {4},
				4: {},
			},
			Ids:         []int{0, 1, 2, 3, 4},
			Entrypoints: []int{0, 3},
			NCycles:     0,
		},
		{
			Name: "two clusters, one with a cycle",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {},
				3: {4},
				4: {3},
			},
			Ids:         []int{0, 1, 2, 3, 4},
			Entrypoints: []int{0},
			NCycles:     1,
		},
		{
			Name: "two clusters, two with a cycle",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {0},
				3: {4},
				4: {3},
			},
			Ids:         []int{0, 1, 2, 3, 4},
			Entrypoints: []int{},
			NCycles:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			testParser := TestParser{
				Spec: tt.Spec,
			}
			var ids []string
			for _, id := range tt.Ids {
				ids = append(ids, strconv.Itoa(id))
			}

			g := NewGraph[[]int]()
			err := g.Load(ids, &testParser, nil)
			a.NoError(err)

			nodesWithoutParents := g.GetNodesWithoutParents()
			cycles := g.RemoveCycles(nodesWithoutParents)

			entrypoints := make([]int, 0)
			for _, entrypoint := range nodesWithoutParents {
				id, _ := strconv.Atoi(entrypoint.Id)
				entrypoints = append(entrypoints, id)
			}
			a.Equal(tt.Entrypoints, entrypoints)
			a.Equal(tt.NCycles, len(cycles))
		})
	}
}
