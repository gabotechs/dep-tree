package graph

import (
	"fmt"
	"os"

	"github.com/gammazero/deque"
	"github.com/schollz/progressbar/v3"
)

type NodeParser[T any] interface {
	Node(id string) (*Node[T], error)
	Deps(node *Node[T]) ([]*Node[T], error)
}

type NodeParserBuilder[T any] func([]string) (NodeParser[T], error)

func (g *Graph[T]) Load(ids []string, parser NodeParser[T], callbacks LoadCallbacks[T]) error {
	if callbacks == nil {
		callbacks = &EmptyCallbacks[T]{}
	}
	visited := make(map[string]bool)
	callbacks.onStartLoading(ids)

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
			visited[node.Id] = true

			deps, err := parser.Deps(node)
			if err != nil {
				node.AddErrors(err)
				continue
			}
			callbacks.onNodeLoaded(node, deps)

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
	onStartLoading(initialIds []string)
	onNodeLoaded(node *Node[T], deps []*Node[T])
	onFinishLoad()
}

type EmptyCallbacks[T any] struct{}

func (e EmptyCallbacks[T]) onStartLoading([]string)               {}
func (e EmptyCallbacks[T]) onNodeLoaded(_ *Node[T], _ []*Node[T]) {}
func (e EmptyCallbacks[T]) onFinishLoad()                         {}

type TestCallbacks[T any] struct {
	startLoad  int
	finishLoad int
	nodeLoaded int
}

func (t *TestCallbacks[T]) onStartLoading([]string)               { t.startLoad++ }
func (t *TestCallbacks[T]) onNodeLoaded(_ *Node[T], _ []*Node[T]) { t.nodeLoaded++ }
func (t *TestCallbacks[T]) onFinishLoad()                         { t.finishLoad++ }

type StdErrCallbacks[T any] struct {
	bar     *progressbar.ProgressBar
	seen    map[string]struct{}
	done    map[string]struct{}
	display func(node *Node[T]) string
}

func NewStdErrCallbacks[T any](
	display func(node *Node[T]) string,
) *StdErrCallbacks[T] {
	bar := progressbar.NewOptions(
		-1,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowIts(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetRenderBlankState(true),
	)
	return &StdErrCallbacks[T]{
		bar,
		map[string]struct{}{},
		map[string]struct{}{},
		display,
	}
}

func (s *StdErrCallbacks[T]) onStartLoading(initialNodes []string) {
	s.bar.Reset()
	for _, nodeId := range initialNodes {
		s.seen[nodeId] = struct{}{}
	}
}

func (s *StdErrCallbacks[T]) onNodeLoaded(n *Node[T], deps []*Node[T]) {
	_ = s.bar.Add(1)
	s.done[n.Id] = struct{}{}
	for _, n := range deps {
		s.seen[n.Id] = struct{}{}
	}
	s.bar.Describe(fmt.Sprintf("(%d/%d) Loading Files...", len(s.done), len(s.seen)))
}

func (s *StdErrCallbacks[T]) onFinishLoad() {
	s.bar.Describe(fmt.Sprintf("(%d/%d) Finished loading", len(s.done), len(s.seen)))
	_ = s.bar.Finish()
	_, _ = os.Stderr.Write([]byte{'\n'})
}
