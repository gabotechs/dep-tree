package dep_tree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadDeps_noOwnChild(t *testing.T) {
	a := require.New(t)
	testGraph := &TestParser{
		Start: "0",
		Spec:  [][]int{{0}},
	}
	dt := NewDepTree[[]int](testGraph)
	root, err := dt.Root()
	a.NoError(err)
	err = dt.LoadDeps()
	a.NoError(err)

	a.Equal(testGraph.Start, root.Id)
	a.Equal(0, len(dt.Graph.ToId(testGraph.Start)))
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

	dt := NewDepTree[[]int](testGraph)

	err := dt.LoadDeps()
	a.NoError(err)
	node0 := dt.Graph.Get("0")
	a.Equal(len(node0.Errors), 0)
	node1 := dt.Graph.Get("1")
	a.Equal(len(node1.Errors), 0)
	node2 := dt.Graph.Get("2")
	a.ErrorContains(node2.Errors[0], "no negative children")
}
