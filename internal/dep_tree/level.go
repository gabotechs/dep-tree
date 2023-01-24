package dep_tree

import (
	"context"
	"fmt"
	"sort"

	"dep-tree/internal/graph"
)

type cacheKey string
type cycleKey string

const unknown = -2
const cyclic = -1

func copyMap[K comparable, V any](m map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		result[k] = v
	}
	return result
}

func hashDep(a string, b string) string {
	return a + " -> " + b
}

func calculateLevel[T any](
	ctx context.Context,
	g *graph.Graph[T],
	nodeId string,
	rootId string,
	seen map[string]bool,
) (context.Context, int) {
	var cachedLevelKey = cacheKey("level-" + rootId + "-" + nodeId)
	if cachedLevel, ok := ctx.Value(cachedLevelKey).(int); ok {
		// 1. Check first the cache, we do not like to work more than need.
		return ctx, cachedLevel
	} else if nodeId == rootId {
		// 2. If it is the root node where are done.
		return ctx, 0
	} else if _, ok := seen[nodeId]; ok {
		// 3. Check if we have closed a loop.
		return ctx, cyclic
	}

	// 4. Calculate the maximum level for this node ignore deps that where previously seen as cyclical.
	seen = copyMap(seen)
	seen[nodeId] = true
	maxLevel := unknown
	for _, parent := range g.Parents(nodeId) {
		dep := hashDep(parent.Id, nodeId)

		cachedCycleKey := cycleKey("cycle-" + rootId + "-" + dep)
		if _, ok := ctx.Value(cachedCycleKey).(bool); ok {
			continue
		}

		var level int
		ctx, level = calculateLevel(ctx, g, parent.Id, rootId, seen)
		if level == cyclic {
			ctx = context.WithValue(ctx, cachedCycleKey, true)
		} else if level > maxLevel {
			maxLevel = level
		}
	}
	// 5. If ignoring previously seen cyclical deps we are not able
	//  to tell the level, then recalculate without ignoring them.
	if maxLevel == unknown {
		for _, parent := range g.Parents(nodeId) {
			var level int
			ctx, level = calculateLevel(ctx, g, parent.Id, rootId, seen)
			if level > maxLevel {
				maxLevel = level
			}
		}
	}
	if maxLevel >= 0 {
		ctx = context.WithValue(ctx, cachedLevelKey, maxLevel+1)
	}
	return ctx, maxLevel + 1
}

// Level retrieves the longest path until going to "rootId" avoiding cyclical loops.
func level[T any](
	ctx context.Context,
	g *graph.Graph[T],
	nodeId string,
	rootId string,
) (context.Context, int) {
	ctx, lvl := calculateLevel(ctx, g, nodeId, rootId, map[string]bool{})
	if lvl == unknown {
		// TODO: there is a bug here, there are cases where this is reached.
		msg := "This should not be reachable"
		msg += fmt.Sprintf("\nhappened while calculating level for node %s", nodeId)
		msg += fmt.Sprintf("\nthis node has %d parents", len(g.Parents(nodeId)))

		panic(msg)
	}
	return ctx, lvl
}

func GetDepTreeNodes[T any](
	ctx context.Context,
	g *graph.Graph[T],
	rootId string,
) (context.Context, []*DepTreeNode[T]) {
	allNodes := g.Nodes()
	result := make([]*DepTreeNode[T], len(allNodes))
	for i, n := range allNodes {
		var lvl int
		ctx, lvl = level(ctx, g, n.Id, rootId)
		result[i] = &DepTreeNode[T]{
			Node: n,
			Lvl:  lvl,
		}
	}

	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Lvl == result[j].Lvl {
			return result[i].Node.Id < result[j].Node.Id
		} else {
			return result[i].Lvl < result[j].Lvl
		}
	})
	return ctx, result
}
