package dep_tree

import (
	"context"
	"fmt"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/schollz/progressbar/v3"

	"github.com/gabotechs/dep-tree/internal/graph"
)

type NodeParserBuilder[T any] func(context.Context, string) (context.Context, NodeParser[T], error)

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
	// Info present on DepTree construction.
	NodeParser[T]
	// Info present just after node processing.
	Graph  *graph.Graph[T]
	Nodes  []*DepTreeNode[T]
	Cycles *orderedmap.OrderedMap[[2]string, graph.Cycle]
	// just some internal cache.
	root *graph.Node[T]
	// callbacks
	onStartLoading func()
	onNodeLoad     func(*graph.Node[T])
	onFinishLoad   func()
}

func NewDepTree[T any](parser NodeParser[T]) *DepTree[T] {
	return (&DepTree[T]{
		NodeParser: parser,
		Nodes:      []*DepTreeNode[T]{},
		Graph:      graph.NewGraph[T](),
		Cycles:     orderedmap.NewOrderedMap[[2]string, graph.Cycle](),
	}).withLoader()
}

func (dt *DepTree[T]) withLoader() *DepTree[T] {
	bar := progressbar.Default(-1)
	dt.onStartLoading = func() {
		bar.Reset()
	}
	dt.onNodeLoad = func(n *graph.Node[T]) {
		_ = bar.Add(1)
		bar.Describe(fmt.Sprintf("Loading %s...", dt.NodeParser.Display(n)))
	}
	dt.onFinishLoad = func() {
		bar.Describe("Finished loading")
		_ = bar.Finish()
		_ = bar.Clear()
	}
	return dt
}

func (dt *DepTree[T]) Root() (*graph.Node[T], error) {
	if dt.root != nil {
		return dt.root, nil
	}
	root, err := dt.NodeParser.Entrypoint()
	if err != nil {
		return nil, err
	}
	dt.root = root
	return dt.root, nil
}
