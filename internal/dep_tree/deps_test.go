package dep_tree

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"dep-tree/internal/graph"
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
