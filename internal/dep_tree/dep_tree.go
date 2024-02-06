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

type DisplayResult struct {
	Name  string
	Group string
}

type NodeParser[T any] interface {
	Display(node *graph.Node[T]) DisplayResult
	Node(id string) (*graph.Node[T], error)
	Deps(node *graph.Node[T]) ([]*graph.Node[T], error)
}

type DepTree[T any] struct {
	// Info present on DepTree construction.
	Ids []string
	NodeParser[T]
	// Info present just after node processing.
	Graph       *graph.Graph[T]
	Entrypoints []*graph.Node[T]
	Cycles      *orderedmap.OrderedMap[[2]string, graph.Cycle]
	// callbacks
	onStartLoading   func()
	onNodeStartLoad  func(*graph.Node[T])
	onNodeFinishLoad func(*graph.Node[T], []*graph.Node[T])
	onFinishLoad     func()
}

func NewDepTree[T any](parser NodeParser[T], ids []string) *DepTree[T] {
	return &DepTree[T]{
		Ids:              ids,
		NodeParser:       parser,
		Graph:            graph.NewGraph[T](),
		Cycles:           orderedmap.NewOrderedMap[[2]string, graph.Cycle](),
		onStartLoading:   func() {},
		onNodeStartLoad:  func(_ *graph.Node[T]) {},
		onNodeFinishLoad: func(_ *graph.Node[T], _ []*graph.Node[T]) {},
		onFinishLoad:     func() {},
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
	done := 0
	dt.onStartLoading = func() {
		bar.Reset()
	}
	dt.onNodeStartLoad = func(n *graph.Node[T]) {
		done += 1
		_ = bar.Set(done)
		bar.Describe(fmt.Sprintf("(%d/%d) Loading %s...", done, len(diff), dt.NodeParser.Display(n).Name))
	}
	dt.onNodeFinishLoad = func(n *graph.Node[T], ns []*graph.Node[T]) {
		for _, n := range ns {
			diff[n.Id] = true
		}
		bar.Describe(fmt.Sprintf("(%d/%d) Loading %s...", done, len(diff), dt.NodeParser.Display(n).Name))
	}
	dt.onFinishLoad = func() {
		bar.Describe("Finished loading")
		_ = bar.Finish()
		_ = bar.Clear()
	}
	return dt
}
