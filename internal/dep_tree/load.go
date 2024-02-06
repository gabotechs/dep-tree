package dep_tree

import (
	"github.com/elliotchance/orderedmap/v2"
	"github.com/gammazero/deque"

	"github.com/gabotechs/dep-tree/internal/graph"
)

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
			dt.onNodeFinishLoad(node, deps)

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
	// These are not as nice, as the rule for determining which cycles are trimmed is more arbitrary.
	for _, cycle := range dt.Graph.RemoveJohnsonCycles() {
		dt.Cycles.Set(cycle.Cause, cycle)
	}
}
