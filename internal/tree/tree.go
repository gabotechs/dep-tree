package tree

import (
	"errors"
	"fmt"
	"sort"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/graph"
)

type NodeWithLevel[T any] struct {
	*graph.Node[T]
	Lvl int
}

type Tree[T any] struct {
	*dep_tree.DepTree[T]
	Nodes []*NodeWithLevel[T]
	// cache
	longestPathCache map[string]int
}

func NewTree[T any](dt *dep_tree.DepTree[T]) (*Tree[T], error) {
	if len(dt.Entrypoints) == 0 {
		return nil, errors.New("this functionality requires that at least 1 entrypoint is provided")
	}
	if len(dt.Entrypoints) > 1 {
		return nil, fmt.Errorf("this functionality requires that only 1 entrypoint is provided, but %d where detected. Consider providing a single entrypoint to your program", len(dt.Entrypoints))
	}
	allNodes := dt.Graph.AllNodes()
	tree := Tree[T]{
		dt,
		make([]*NodeWithLevel[T], len(allNodes)),
		make(map[string]int),
	}
	for i, n := range allNodes {
		lvl, err := tree.longestPath(dt.Entrypoints[0].Id, n.Id, nil)
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

func (t *Tree[T]) LoadNodes() error {
	return nil
}
