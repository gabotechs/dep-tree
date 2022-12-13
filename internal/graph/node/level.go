package node

import (
	"context"
	"fmt"

	"dep-tree/internal/utils"
)

type key int

const (
	cycleKey key = iota
)

const unknown = -2
const cyclic = -1

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
		} else if utils.InArray(dep, knownCycles) {
			continue
		}

		var newLevel int
		ctx, newLevel = calculateLevel(ctx, parent, rootId, level+1, append(stack, node.Id))
		if newLevel == cyclic {
			ctx = context.WithValue(ctx, cycleKey, append(knownCycles, dep))
		} else if newLevel > maxLevel {
			maxLevel = newLevel
		}
	}
	if maxLevel == unknown {
		// TODO: there is a bug here, there are cases where this is reached.
		msg := "This should not be reachable"
		msg += fmt.Sprintf("\nhappened while calculating level for node %s", node.Id)
		msg += fmt.Sprintf("\nthis node has %d parents", node.Parents.Len())

		panic(msg)
	} else {
		return ctx, maxLevel
	}
}

// Level retrieves the longest path until going to "rootId" avoiding cyclical loops.
func (n *Node[T]) Level(ctx context.Context, rootId string) (context.Context, int) {
	return calculateLevel(ctx, n, rootId, 0, []string{})
}
