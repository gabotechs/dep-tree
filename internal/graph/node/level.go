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

func hashDep(a string, b string) string {
	return a + " -> " + b
}

func (f *FamilyRegistry[T]) calculateLevel(
	ctx context.Context,
	nodeId string,
	rootId string,
	seen map[string]bool,
) (context.Context, int) {
	var cachedLevelKey = cacheKey("level-" + nodeId)
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
	for _, parent := range f.Parents(nodeId) {
		dep := hashDep(parent.Id, nodeId)

		cachedCycleKey := cycleKey("cycle-" + dep)
		if _, ok := ctx.Value(cachedCycleKey).(bool); ok {
			continue
		}

		var level int
		ctx, level = f.calculateLevel(ctx, parent.Id, rootId, seen)
		if level == cyclic {
			ctx = context.WithValue(ctx, cachedCycleKey, true)
		} else if level > maxLevel {
			maxLevel = level
		}
	}
	// 5. If ignoring previously seen cyclical deps we are not able
	//  to tell the level, then recalculate without ignoring them.
	if maxLevel == unknown {
		for _, parent := range f.Parents(nodeId) {
			var level int
			ctx, level = f.calculateLevel(ctx, parent.Id, rootId, seen)
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
func (f *FamilyRegistry[T]) Level(ctx context.Context, nodeId string, rootId string) (context.Context, int) {
	ctx, lvl := f.calculateLevel(ctx, nodeId, rootId, map[string]bool{})
	if lvl == unknown {
		// TODO: there is a bug here, there are cases where this is reached.
		msg := "This should not be reachable"
		msg += fmt.Sprintf("\nhappened while calculating level for node %s", nodeId)
		msg += fmt.Sprintf("\nthis node has %d parents", len(f.Parents(nodeId)))

		panic(msg)
	}
	return ctx, lvl
}
