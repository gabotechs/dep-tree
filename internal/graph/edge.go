package graph

import (
	"gonum.org/v1/gonum/graph"
)

type Edge[T any] struct {
	from *Node[T]
	to   *Node[T]
}

func (e *Edge[T]) From() graph.Node {
	return e.from
}

func (e *Edge[T]) To() graph.Node {
	return e.to
}

func (e *Edge[T]) ReversedEdge() graph.Edge {
	return &Edge[T]{
		from: e.to,
		to:   e.from,
	}
}
