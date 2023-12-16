package graph

import (
	"github.com/gabotechs/dep-tree/internal/utils"
)

type Cycle struct {
	Cause [2]string
	Stack []string
}

func (g *Graph[T]) removeCycles(node *Node[T], callstack *utils.CallStack) []Cycle {
	err := callstack.Push(node.Id)
	if err != nil {
		last, _ := callstack.Back()
		g.RemoveFromToEdge(last, node.Id)
		var stack []string
		addFlag := false
		for _, el := range callstack.Stack() {
			if el == node.Id {
				addFlag = true
			}
			if addFlag {
				stack = append(stack, el)
			}
		}

		return []Cycle{{
			Cause: [2]string{last, node.Id},
			Stack: append(stack, node.Id),
		}}
	}
	var cycles []Cycle
	for _, toNode := range g.FromId(node.Id) {
		cycles = append(cycles, g.removeCycles(toNode, callstack)...)
	}
	callstack.Pop()
	return cycles
}

func (g *Graph[T]) RemoveCycles(node *Node[T]) []Cycle {
	return g.removeCycles(node, utils.NewCallStack())
}
