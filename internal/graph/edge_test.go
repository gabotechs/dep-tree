package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/graph"
)

func TestEdge_ReversedEdge(t *testing.T) {
	a := require.New(t)
	var edge graph.Edge = &Edge[int]{
		from: MakeNode("1", 1),
		to:   MakeNode("2", 2),
	}

	edge = edge.ReversedEdge()

	a.Equal(edge.To().ID(), hash("1"))
	a.Equal(edge.From().ID(), hash("2"))
}
