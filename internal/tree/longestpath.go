package tree

import (
	"errors"

	"github.com/gabotechs/dep-tree/internal/utils"
)

// longestPath finds the longest path between two nodes in a directed acyclic graph.
//
//	It uses the context as cache, so next calculations even with different nodes are quick.
func (t *Tree[T]) longestPath(
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
	if cachedLevel, ok := t.longestPathCache[cachedLevelKey]; ok {
		return cachedLevel, nil
	}
	err := stack.Push(nodeId)
	if err != nil {
		return 0, errors.New("cannot calculate longest path between nodes because there is at least one cycle in the graph: " + err.Error())
	}

	maxLongestPath := 0
	for _, from := range t.Graph.ToId(nodeId) {
		length, err := t.longestPath(rootId, from.Id, stack)
		if err != nil {
			return 0, err
		}
		if length > maxLongestPath {
			maxLongestPath = length
		}
	}
	if maxLongestPath >= 0 {
		t.longestPathCache[cachedLevelKey] = maxLongestPath + 1
	}

	stack.Pop()
	return maxLongestPath + 1, nil
}
