package tree

import (
	"errors"
	"fmt"
	"sort"

	"github.com/gabotechs/dep-tree/internal/graph"
)

type NodeWithLevel[T any] struct {
	*graph.Node[T]
	Lvl int
}

type Tree[T any] struct {
	Graph      *graph.Graph[T]
	NodeParser graph.NodeParser[T]

	display          func(node *graph.Node[T]) string
	entrypoint       *graph.Node[T]
	Nodes            []*NodeWithLevel[T]
	Cycles           []graph.Cycle
	longestPathCache map[string]int
}

func NewTree[T any](
	files []string,
	parser graph.NodeParser[T],
	display func(node *graph.Node[T]) string,
	callbacks graph.LoadCallbacks[T],
) (*Tree[T], error) {
	if len(files) == 0 {
		return nil, errors.New("this functionality requires that at least 1 entrypoint is provided")
	}
	if len(files) > 1 {
		return nil, fmt.Errorf("this functionality requires that only 1 entrypoint is provided, but %d where passed. Consider providing a single entrypoint to your program", len(files))
	}
	g := graph.NewGraph[T]()
	err := g.Load(files, parser, callbacks)
	if err != nil {
		return nil, err
	}
	entrypoint, err := parser.Node(files[0])
	if err != nil {
		return nil, err
	}

	cycles := g.RemoveCyclesStartingFromNode(entrypoint)

	allNodes := g.AllNodes()
	tree := Tree[T]{
		Graph:            g,
		NodeParser:       parser,
		display:          display,
		entrypoint:       entrypoint,
		Nodes:            make([]*NodeWithLevel[T], len(allNodes)),
		Cycles:           cycles,
		longestPathCache: make(map[string]int),
	}
	for i, n := range allNodes {
		lvl, err := tree.longestPath(entrypoint.Id, n.Id, nil)
		if err != nil {
			return nil, err
		}
		tree.Nodes[i] = &NodeWithLevel[T]{n, lvl}
	}

	sort.SliceStable(tree.Nodes, func(i, j int) bool {
		if tree.Nodes[i].Lvl == tree.Nodes[j].Lvl {
			return tree.Nodes[i].Node.Id < tree.Nodes[j].Node.Id
		} else {
			return tree.Nodes[i].Lvl < tree.Nodes[j].Lvl
		}
	})
	return &tree, nil
}
