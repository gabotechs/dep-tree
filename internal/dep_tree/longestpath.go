package dep_tree

import (
	"context"
	"errors"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/utils"
)

type cacheKey string

// longestPath finds the longest path between two nodes in a directed acyclic graph.
//
//	It uses the context as cache, so next calculations even with different nodes are quick.
func longestPath[T any](
	ctx context.Context,
	g *graph.Graph[T],
	rootId string,
	nodeId string,
	stack *utils.CallStack,
) (context.Context, int, error) {
	if stack == nil {
		stack = utils.NewCallStack()
	}
	if nodeId == rootId {
		return ctx, 0, nil
	}
	var cachedLevelKey = cacheKey(rootId + "-" + nodeId)
	if cachedLevel, ok := ctx.Value(cachedLevelKey).(int); ok {
		return ctx, cachedLevel, nil
	}
	err := stack.Push(nodeId)
	if err != nil {
		return nil, 0, errors.New("cannot calculate longest path between nodes because there is at least one cycle in the graph: " + err.Error())
	}

	maxLongestPath := 0
	for _, from := range g.ToId(nodeId) {
		newCtx, length, err := longestPath(ctx, g, rootId, from.Id, stack)
		ctx = newCtx
		if err != nil {
			return ctx, 0, err
		}
		if length > maxLongestPath {
			maxLongestPath = length
		}
	}
	if maxLongestPath >= 0 {
		ctx = context.WithValue(ctx, cachedLevelKey, maxLongestPath+1)
	}

	stack.Pop()
	return ctx, maxLongestPath + 1, nil
}
