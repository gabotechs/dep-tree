package dep_tree

import (
	"errors"
	"fmt"
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

	visited := make(map[string]bool)
	dt.onStartLoading()

	for _, id := range dt.Ids {
		node, err := dt.NodeParser.Node(id)
		if err != nil {
			return err
		}
		var queue deque.Deque[*graph.Node[T]]
		queue.PushBack(node)
		if !dt.Graph.Has(node.Id) {
			dt.Graph.AddNode(node)
		}
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
	}
	if len(dt.Ids) == 1 {
		// If exactly one file was provided, take that as the entrypoint.
		dt.Entrypoints = []*graph.Node[T]{dt.Graph.Get(dt.Ids[0])}
	} else {
		// If multiple files were provided, use the nodes without parents as entrypoints.
		// Note that due to cyclic dependencies, there might be no parents without dependencies.
		dt.Entrypoints = dt.Graph.GetNodesWithoutParents()
	}
	dt.onFinishLoad()

	return nil
}

func (dt *DepTree[T]) LoadCycles() {
	dt.Cycles = orderedmap.NewOrderedMap[[2]string, graph.Cycle]()

	// First, remove the cycles computed from each entrypoint. This allows
	// us trim the cycles in a more "controlled way"
	for _, entrypoint := range dt.Entrypoints {
		for _, cycle := range dt.Graph.RemoveCycles(entrypoint) {
			dt.Cycles.Set(cycle.Cause, cycle)
		}
	}
	// Then, remove the cycles computed without taking entrypoints into account.
	// These are not as nice as the rule for which cycles are trimmed is more arbitrary.
	for _, cycle := range dt.Graph.RemoveJohnsonCycles() {
		dt.Cycles.Set(cycle.Cause, cycle)
	}
}

func (dt *DepTree[T]) LoadNodes() error {
	if len(dt.Entrypoints) == 0 {
		return errors.New("this functionality requires that at least 1 entrypoint is provided")
	}
	if len(dt.Entrypoints) > 1 {
		return fmt.Errorf("this functionality requires that only 1 entrypoint is provided, but %d where detected. Consider providing a single entrypoint to your program", len(dt.Entrypoints))
	}
	allNodes := dt.Graph.AllNodes()
	dt.Nodes = make([]*DepTreeNode[T], len(allNodes))
	for i, n := range allNodes {
		lvl, err := dt.longestPath(dt.Graph, dt.Entrypoints[0].Id, n.Id, nil)
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
