package dep_tree

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/graph"
)

func TestLoadDeps_noOwnChild(t *testing.T) {
	a := require.New(t)
	testGraph := &TestParser{
		Start: "0",
		Spec:  [][]int{{0}},
	}

	g := graph.NewGraph[[]int]()

	_, rootId, err := LoadDeps[[]int](context.Background(), g, testGraph)
	a.NoError(err)

	a.Equal(testGraph.Start, rootId)
	a.Equal(0, len(g.Children(testGraph.Start)))
}

func TestLoadDeps_ErrorHandle(t *testing.T) {
	a := require.New(t)
	testGraph := &TestParser{
		Start: "0",
		Spec: [][]int{
			{1},
			{2},
			{-3},
		},
	}

	g := graph.NewGraph[[]int]()

	_, _, err := LoadDeps[[]int](context.Background(), g, testGraph)
	a.NoError(err)
	node0 := g.Get("0")
	a.Equal(len(node0.Errors), 0)
	node1 := g.Get("1")
	a.Equal(len(node1.Errors), 0)
	node2 := g.Get("2")
	a.ErrorContains(node2.Errors[0], "no negative children")
}
