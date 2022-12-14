package node

import (
	"context"
	"fmt"
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

func hashDep[T any](a *Node[T], b *Node[T]) string {
	return a.Id + " -> " + b.Id
}

func calculateLevel[T any](
	ctx context.Context,
	node *Node[T],
	rootId string,
	level int,
	seen map[string]bool,
) (context.Context, int) {
	var cachedLevelKey = cacheKey("level-" + node.Id)
	if cachedLevel, ok := ctx.Value(cachedLevelKey).(int); ok {
		// 1. Check first the cache, we do not like to work more than need.
		return ctx, cachedLevel + level
	} else if node.Id == rootId {
		// 2. If it is the root node where are done.
		return ctx, level
	} else if _, ok := seen[node.Id]; ok {
		// 3. Check if we have closed a loop.
		return ctx, cyclic
	}

	// 4. Calculate the maximum level for this node ignore deps that where previously seen as cyclical.
	seen = copyMap(seen)
	seen[node.Id] = true
	maxLevel := unknown
	for _, parentId := range node.Parents.Keys() {
		parent, _ := node.Parents.Get(parentId)
		dep := hashDep(parent, node)

		cachedCycleKey := cycleKey("cycle-" + dep)
		if _, ok := ctx.Value(cachedCycleKey).(bool); ok {
			continue
		}

		var newLevel int
		ctx, newLevel = calculateLevel(ctx, parent, rootId, level+1, seen)
		if newLevel == cyclic {
			ctx = context.WithValue(ctx, cachedCycleKey, true)
		} else if newLevel > maxLevel {
			maxLevel = newLevel
		}
	}
	// 5. If ignoring previously seen cyclical deps we are not able
	//  to tell the level, then recalculate without ignoring them.
	if maxLevel == unknown {
		for _, parentId := range node.Parents.Keys() {
			parent, _ := node.Parents.Get(parentId)

			var newLevel int
			ctx, newLevel = calculateLevel(ctx, parent, rootId, level+1, seen)
			if newLevel > maxLevel {
				maxLevel = newLevel
			}
		}
	}
	//if maxLevel >= 0 {
	//	ctx = context.WithValue(ctx, cachedLevelKey, maxLevel)
	//}
	return ctx, maxLevel
}

// Level retrieves the longest path until going to "rootId" avoiding cyclical loops.
func (n *Node[T]) Level(ctx context.Context, rootId string) (context.Context, int) {
	ctx, lvl := calculateLevel(ctx, n, rootId, 0, map[string]bool{})
	if lvl == unknown {
		// TODO: there is a bug here, there are cases where this is reached.
		msg := "This should not be reachable"
		msg += fmt.Sprintf("\nhappened while calculating level for node %s", n.Id)
		msg += fmt.Sprintf("\nthis node has %d parents", n.Parents.Len())

		panic(msg)
	}
	return ctx, lvl
}
