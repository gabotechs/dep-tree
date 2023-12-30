package dep_tree

import (
	"sort"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/gammazero/deque"

	"github.com/gabotechs/dep-tree/internal/graph"
)

func (dt *DepTree[T]) LoadDeps() error {
	err := dt.LoadGraph()
	if err != nil {
		return err
	}

	dt.LoadCycles()

	return dt.LoadNodes()
}

func (dt *DepTree[T]) LoadGraph() error {
	dt.Graph = graph.NewGraph[T]()

	root, err := dt.Root()
	if err != nil {
		return err
	}

	var queue deque.Deque[*graph.Node[T]]
	queue.PushBack(root)
	dt.Graph.AddNode(root)

	visited := make(map[string]bool)
	dt.onStartLoading()

	for queue.Len() > 0 {
		node := queue.PopFront()
		if _, ok := visited[node.Id]; ok {
			continue
		}
		dt.onNodeStartLoad(node)
		visited[node.Id] = true

		deps, err := dt.NodeParser.Deps(node)
		if err != nil {
			node.AddErrors(err)
			continue
		}
		dt.onNodeFinishLoad(deps)

		for _, dep := range deps {
			// No own child.
			if dep.Id == node.Id {
				continue
			}
			if !dt.Graph.Has(dep.Id) {
				dt.Graph.AddNode(dep)
			}
			err = dt.Graph.AddFromToEdge(node.Id, dep.Id)
			queue.PushBack(dep)
			if err != nil {
				return err
			}
		}
	}
	dt.onFinishLoad()

	return nil
}

func (dt *DepTree[T]) LoadCycles() {
	dt.Cycles = orderedmap.NewOrderedMap[[2]string, graph.Cycle]()

	cycles := dt.Graph.RemoveCycles(dt.root)
	for _, cycle := range cycles {
		dt.Cycles.Set(cycle.Cause, cycle)
	}
}

func (dt *DepTree[T]) LoadNodes() error {
	root, err := dt.Root()
	if err != nil {
		return err
	}
	allNodes := dt.Graph.AllNodes()
	dt.Nodes = make([]*DepTreeNode[T], len(allNodes))
	for i, n := range allNodes {
		lvl, err := dt.longestPath(dt.Graph, root.Id, n.Id, nil)
		if err != nil {
			return err
		}
		dt.Nodes[i] = &DepTreeNode[T]{n, lvl}
	}

	sort.SliceStable(dt.Nodes, func(i, j int) bool {
		if dt.Nodes[i].Lvl == dt.Nodes[j].Lvl {
			return dt.Nodes[i].Node.Id < dt.Nodes[j].Node.Id
		} else {
			return dt.Nodes[i].Lvl < dt.Nodes[j].Lvl
		}
	})
	return nil
}
