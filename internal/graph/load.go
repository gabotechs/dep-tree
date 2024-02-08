package graph

import (
	"fmt"
	"os"
	"time"

	"github.com/gammazero/deque"
	"github.com/schollz/progressbar/v3"
)

type DisplayResult struct {
	Name  string
	Group string
}

type NodeParser[T any] interface {
	Display(node *Node[T]) DisplayResult
	Node(id string) (*Node[T], error)
	Deps(node *Node[T]) ([]*Node[T], error)
}

type NodeParserBuilder[T any] func([]string) (NodeParser[T], error)

func (g *Graph[T]) Load(ids []string, parser NodeParser[T], callbacks LoadCallbacks[T]) error {
	if callbacks == nil {
		callbacks = &EmptyCallbacks[T]{}
	}
	visited := make(map[string]bool)
	callbacks.onStartLoading()

	for _, id := range ids {
		node, err := parser.Node(id)
		if err != nil {
			return err
		}
		var queue deque.Deque[*Node[T]]
		queue.PushBack(node)
		if !g.Has(node.Id) {
			g.AddNode(node)
		}
		for queue.Len() > 0 {
			node := queue.PopFront()
			if _, ok := visited[node.Id]; ok {
				continue
			}
			callbacks.onNodeStartLoad(node)
			visited[node.Id] = true

			deps, err := parser.Deps(node)
			if err != nil {
				node.AddErrors(err)
				continue
			}
			callbacks.onNodeFinishLoad(node, deps)

			for _, dep := range deps {
				// No own child.
				if dep.Id == node.Id {
					continue
				}
				if !g.Has(dep.Id) {
					g.AddNode(dep)
				}
				err = g.AddFromToEdge(node.Id, dep.Id)
				queue.PushBack(dep)
				if err != nil {
					return err
				}
			}
		}
	}
	callbacks.onFinishLoad()

	return nil
}

type LoadCallbacks[T any] interface {
	onStartLoading()
	onNodeStartLoad(node *Node[T])
	onNodeFinishLoad(node *Node[T], deps []*Node[T])
	onFinishLoad()
}

type EmptyCallbacks[T any] struct{}

func (e EmptyCallbacks[T]) onStartLoading()                           {}
func (e EmptyCallbacks[T]) onNodeStartLoad(_ *Node[T])                {}
func (e EmptyCallbacks[T]) onNodeFinishLoad(_ *Node[T], _ []*Node[T]) {}
func (e EmptyCallbacks[T]) onFinishLoad()                             {}

type TestCallbacks[T any] struct {
	startLoad  int
	startNode  int
	finishLoad int
	finishNode int
}

func (t *TestCallbacks[T]) onStartLoading() {
	t.startLoad++
}
func (t *TestCallbacks[T]) onNodeStartLoad(_ *Node[T]) {
	t.startNode++
}
func (t *TestCallbacks[T]) onNodeFinishLoad(_ *Node[T], _ []*Node[T]) {
	t.finishNode++
}
func (t *TestCallbacks[T]) onFinishLoad() {
	t.finishLoad++
}

type StdErrCallbacks[T any] struct {
	bar   *progressbar.ProgressBar
	nodes map[string]struct{}
	done  int
}

func NewStdErrCallbacks[T any]() *StdErrCallbacks[T] {
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
	return &StdErrCallbacks[T]{bar, map[string]struct{}{}, 0}
}

func (s *StdErrCallbacks[T]) onStartLoading() {
	s.bar.Reset()
	s.nodes = map[string]struct{}{}
	s.done = 0
}
func (s *StdErrCallbacks[T]) onNodeStartLoad(node *Node[T]) {
	s.done += 1
	_ = s.bar.Set(s.done)
	s.bar.Describe(fmt.Sprintf("(%d/%d) Loading %s...", s.done, len(s.nodes), node.Id))
}
func (s *StdErrCallbacks[T]) onNodeFinishLoad(node *Node[T], deps []*Node[T]) {
	for _, n := range deps {
		s.nodes[n.Id] = struct{}{}
	}
	s.bar.Describe(fmt.Sprintf("(%d/%d) Loading %s...", s.done, len(s.nodes), node.Id))
}
func (s *StdErrCallbacks[T]) onFinishLoad() {
	s.bar.Describe("Finished loading")
	_ = s.bar.Finish()
	_ = s.bar.Clear()
}
