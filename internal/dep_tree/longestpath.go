package dep_tree

import (
	"errors"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/utils"
)

// longestPath finds the longest path between two nodes in a directed acyclic graph.
//
//	It uses the context as cache, so next calculations even with different nodes are quick.
func (dt *DepTree[T]) longestPath(
	g *graph.Graph[T],
	rootId string,
	nodeId string,
	stack *utils.CallStack,
) (int, error) {
	if stack == nil {
		stack = utils.NewCallStack()
	}
	if nodeId == rootId {
		return 0, nil
	}
	var cachedLevelKey = rootId + "-" + nodeId
	if cachedLevel, ok := dt.longestPathCache[cachedLevelKey]; ok {
		return cachedLevel, nil
	}
	err := stack.Push(nodeId)
	if err != nil {
		return 0, errors.New("cannot calculate longest path between nodes because there is at least one cycle in the graph: " + err.Error())
	}

	maxLongestPath := 0
	for _, from := range g.ToId(nodeId) {
		length, err := dt.longestPath(g, rootId, from.Id, stack)
		if err != nil {
			return 0, err
		}
		if length > maxLongestPath {
			maxLongestPath = length
		}
	}
	if maxLongestPath >= 0 {
		dt.longestPathCache[cachedLevelKey] = maxLongestPath + 1
	}

	stack.Pop()
	return maxLongestPath + 1, nil
}
