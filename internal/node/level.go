package node

import "context"

const cycleKey = "cycle"

const unknown = -2
const cyclic = -1

func elementInArray[T comparable](value T, array []T) bool {
	for _, el := range array {
		if value == el {
			return true
		}
	}
	return false
}

func hashDep[T any](a *Node[T], b *Node[T]) string {
	return a.Id + " -> " + b.Id
}

func calculateLevel[T any](
	ctx context.Context,
	node *Node[T],
	rootId string,
	level int,
	stack []string,
) (context.Context, int) {
	if node.Id == rootId {
		return ctx, level
	} else {
		for _, seen := range stack {
			if seen == node.Id {
				return ctx, cyclic
			}
		}
	}
	maxLevel := unknown
	for _, parentId := range node.Parents.Keys() {
		parent, _ := node.Parents.Get(parentId)
		dep := hashDep(parent, node)
		knownCycles, _ := ctx.Value(cycleKey).([]string)
		if knownCycles == nil {
			knownCycles = []string{}
		} else if elementInArray(dep, knownCycles) {
			continue
		}

		var newLevel int
		ctx, newLevel = calculateLevel(ctx, parent, rootId, level+1, append(stack, node.Id))
		if newLevel == cyclic {
			ctx = context.WithValue(ctx, cycleKey, append(knownCycles, dep))
			continue
		} else if newLevel > maxLevel {
			maxLevel = newLevel
		}
	}
	if maxLevel == unknown {
		panic("This should not be reachable")
	} else {
		return ctx, maxLevel
	}
}

// Level retrieves the longest path until going to "rootId" avoiding cyclical loops
func (n *Node[T]) Level(ctx context.Context, rootId string) (context.Context, int) {
	return calculateLevel(ctx, n, rootId, 0, []string{})
}
