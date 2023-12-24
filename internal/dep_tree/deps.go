package dep_tree

import (
	"context"
	"sort"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/gammazero/deque"

	"github.com/gabotechs/dep-tree/internal/graph"
)

func (dt *DepTree[T]) LoadDeps(ctx context.Context) (context.Context, error) {
	ctx, err := dt.LoadGraph(ctx)
	if err != nil {
		return ctx, err
	}

	dt.LoadCycles()

	return dt.LoadNodes(ctx)
}

func (dt *DepTree[T]) LoadGraph(ctx context.Context) (context.Context, error) {
	dt.Graph = graph.NewGraph[T]()

	root, err := dt.Root()
	if err != nil {
		return ctx, err
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

		newCtx, deps, err := dt.NodeParser.Deps(ctx, node)
		ctx = newCtx
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
				return ctx, err
			}
		}
	}
	dt.onFinishLoad()

	return ctx, nil
}

func (dt *DepTree[T]) LoadCycles() {
	dt.Cycles = orderedmap.NewOrderedMap[[2]string, graph.Cycle]()

	cycles := dt.Graph.RemoveCycles(dt.root)
	for _, cycle := range cycles {
		dt.Cycles.Set(cycle.Cause, cycle)
	}
}

func (dt *DepTree[T]) LoadNodes(ctx context.Context) (context.Context, error) {
	root, err := dt.Root()
	if err != nil {
		return ctx, err
	}
	allNodes := dt.Graph.AllNodes()
	dt.Nodes = make([]*DepTreeNode[T], len(allNodes))
	for i, n := range allNodes {
		newCtx, lvl, err := longestPath(ctx, dt.Graph, root.Id, n.Id, nil)
		ctx = newCtx
		if err != nil {
			return ctx, err
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
	return ctx, nil
}
