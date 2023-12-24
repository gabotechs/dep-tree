package dep_tree

import (
	"context"
	"fmt"
	"os"
	"time"

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
	onStartLoading   func()
	onNodeStartLoad  func(*graph.Node[T])
	onNodeFinishLoad func([]*graph.Node[T])
	onFinishLoad     func()
}

func NewDepTree[T any](parser NodeParser[T]) *DepTree[T] {
	return (&DepTree[T]{
		NodeParser:       parser,
		Nodes:            []*DepTreeNode[T]{},
		Graph:            graph.NewGraph[T](),
		Cycles:           orderedmap.NewOrderedMap[[2]string, graph.Cycle](),
		onStartLoading:   func() {},
		onNodeStartLoad:  func(_ *graph.Node[T]) {},
		onNodeFinishLoad: func(_ []*graph.Node[T]) {},
		onFinishLoad:     func() {},
	}).withLoader()
}

func (dt *DepTree[T]) withLoader() *DepTree[T] {
	bar := progressbar.NewOptions64(
		-1,
		progressbar.OptionSetDescription("Loading graph..."),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowIts(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)
	diff := make(map[string]bool)
	total := 0
	dt.onStartLoading = func() {
		bar.Reset()
	}
	dt.onNodeStartLoad = func(n *graph.Node[T]) {
		total += 1
		_ = bar.Set(total)
		bar.Describe(fmt.Sprintf("(%d/%d) Loading %s...", total, len(diff), dt.NodeParser.Display(n)))
	}
	dt.onNodeFinishLoad = func(ns []*graph.Node[T]) {
		for _, n := range ns {
			diff[n.Id] = true
		}
		bar.ChangeMax(len(diff))
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
