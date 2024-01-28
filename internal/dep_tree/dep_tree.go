package dep_tree

import (
	"fmt"
	"os"
	"time"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/schollz/progressbar/v3"

	"github.com/gabotechs/dep-tree/internal/graph"
)

type NodeParserBuilder[T any] func([]string) (NodeParser[T], error)

type NodeParser[T any] interface {
	Display(node *graph.Node[T]) string
	Node(id string) (*graph.Node[T], error)
	Deps(node *graph.Node[T]) ([]*graph.Node[T], error)
}

type DepTreeNode[T any] struct {
	Node *graph.Node[T]
	Lvl  int
}

type DepTree[T any] struct {
	// Info present on DepTree construction.
	Ids []string
	NodeParser[T]
	// Info present just after node processing.
	Graph       *graph.Graph[T]
	Entrypoints []*graph.Node[T]
	Nodes       []*DepTreeNode[T]
	Cycles      *orderedmap.OrderedMap[[2]string, graph.Cycle]
	// callbacks
	onStartLoading   func()
	onNodeStartLoad  func(*graph.Node[T])
	onNodeFinishLoad func([]*graph.Node[T])
	onFinishLoad     func()
	// cache
	longestPathCache map[string]int
}

func NewDepTree[T any](parser NodeParser[T], ids []string) *DepTree[T] {
	return &DepTree[T]{
		Ids:              ids,
		NodeParser:       parser,
		Nodes:            []*DepTreeNode[T]{},
		Graph:            graph.NewGraph[T](),
		Cycles:           orderedmap.NewOrderedMap[[2]string, graph.Cycle](),
		onStartLoading:   func() {},
		onNodeStartLoad:  func(_ *graph.Node[T]) {},
		onNodeFinishLoad: func(_ []*graph.Node[T]) {},
		onFinishLoad:     func() {},
		longestPathCache: map[string]int{},
	}
}

func (dt *DepTree[T]) WithStdErrLoader() *DepTree[T] {
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
