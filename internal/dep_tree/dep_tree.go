package dep_tree

import (
	"context"

	"dep-tree/internal/graph"
)

type NodeParser[T any] interface {
	Display(node *graph.Node[T]) string
	Entrypoint() (*graph.Node[T], error)
	Deps(ctx context.Context, node *graph.Node[T]) (context.Context, []*graph.Node[T], error)
}

type DepTreeNode[T any] struct {
	Node *graph.Node[T]
	Lvl  int
}

type DepTree[T any] struct {
	Nodes  []*DepTreeNode[T]
	Graph  *graph.Graph[T]
	RootId string
	Cycles [][]string
}

func NewDepTree[T any](
	ctx context.Context,
	parser NodeParser[T],
) (context.Context, *DepTree[T], error) {
	// 1. create graph.
	g := graph.NewGraph[T]()
	// 2. populate the graph.
	ctx, rootId, err := LoadDeps(ctx, g, parser)
	if err != nil {
		return ctx, nil, err
	}
	// 3. get sorted by level.
	ctx, nodes := GetDepTreeNodes(ctx, g, rootId)
	return ctx, &DepTree[T]{
		Nodes:  nodes,
		Graph:  g,
		RootId: rootId,
		Cycles: [][]string{}, // TODO.
	}, nil
}
