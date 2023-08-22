package dep_tree

import (
	"context"
	"fmt"
	"sort"

	"github.com/elliotchance/orderedmap/v2"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/utils"
)

type cacheKey string

const unknown = -2
const cyclic = -1

type DepCycle struct {
	Cause [2]string
	Stack []string
}

type LevelCalculator[T any] struct {
	g      *graph.Graph[T]
	rootId string
	Cycles *orderedmap.OrderedMap[[2]string, DepCycle]
}

func NewLevelCalculator[T any](
	g *graph.Graph[T],
	rootId string,
) *LevelCalculator[T] {
	return &LevelCalculator[T]{
		g:      g,
		rootId: rootId,
		Cycles: orderedmap.NewOrderedMap[[2]string, DepCycle](),
	}
}

func (lc *LevelCalculator[T]) calculateLevel(
	ctx context.Context,
	nodeId string,
	stack []string,
) (context.Context, int) {
	var cachedLevelKey = cacheKey("level-" + lc.rootId + "-" + nodeId)
	if nodeId == lc.rootId {
		// 1. If it is the root node where are done.
		return ctx, 0
	} else if cachedLevel, ok := ctx.Value(cachedLevelKey).(int); ok {
		// 2. Check first the cache, we do not like to work more than need.
		return ctx, cachedLevel
	} else if utils.InArray(nodeId, stack) {
		// 3. Check if we have closed a loop.
		return ctx, cyclic
	}
	stack = append([]string{nodeId}, stack...) // reverse because we go from child to parent.

	// 4. Calculate the maximum level for this node ignoring deps that where previously seen as cyclical.
	maxLevel := unknown
	for _, parent := range lc.g.Parents(nodeId) {
		dep := [2]string{parent.Id, nodeId}
		if _, ok := lc.Cycles.Get(dep); ok {
			continue
		}
		var level int
		ctx, level = lc.calculateLevel(ctx, parent.Id, stack)
		if level == cyclic {
			cycleStack := []string{parent.Id}
			for _, stackElement := range stack {
				cycleStack = append(cycleStack, stackElement)
				if stackElement == parent.Id {
					break
				}
			}

			lc.Cycles.Set(dep, DepCycle{
				Cause: dep,
				Stack: cycleStack,
			})
		} else if level > maxLevel {
			maxLevel = level
		}
	}
	// 5. If ignoring previously seen cyclical deps we are not able
	//  to tell the level, then recalculate without ignoring them.
	if maxLevel == unknown {
		for _, parent := range lc.g.Parents(nodeId) {
			var level int
			ctx, level = lc.calculateLevel(ctx, parent.Id, stack)
			if level > maxLevel {
				maxLevel = level
			}
		}
	}
	if maxLevel >= 0 {
		ctx = context.WithValue(ctx, cachedLevelKey, maxLevel+1)
	} else if maxLevel == unknown {
		return ctx, unknown
	}
	return ctx, maxLevel + 1
}

// Level retrieves the longest path until going to "rootId" avoiding cyclical loops.
func (lc *LevelCalculator[T]) level(
	ctx context.Context,
	nodeId string,
) (context.Context, int) {
	ctx, lvl := lc.calculateLevel(ctx, nodeId, []string{})
	if lvl == unknown {
		// TODO: there is a bug here, there are cases where this is reached.
		msg := "This should not be reachable"
		msg += fmt.Sprintf("\nhappened while calculating level for node %s", nodeId)
		msg += fmt.Sprintf("\nthis node has %d parents", len(lc.g.Parents(nodeId)))
		panic(msg)
	}
	return ctx, lvl
}

func GetDepTreeNodes[T any](
	ctx context.Context,
	g *graph.Graph[T],
	rootId string,
) (context.Context, []*DepTreeNode[T], *orderedmap.OrderedMap[[2]string, DepCycle]) {
	levelCalculator := NewLevelCalculator(g, rootId)

	allNodes := g.Nodes()
	result := make([]*DepTreeNode[T], len(allNodes))
	for i, n := range allNodes {
		var lvl int
		ctx, lvl = levelCalculator.level(ctx, n.Id)
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

	return ctx, result, levelCalculator.Cycles
}
