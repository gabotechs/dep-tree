package graph

import (
	"gonum.org/v1/gonum/graph/topo"

	"github.com/gabotechs/dep-tree/internal/utils"
)

type Cycle struct {
	Cause [2]string
	Stack []string
}

func (g *Graph[T]) removeCycles(node *Node[T], callstack *utils.CallStack, done map[string]bool) []Cycle {
	if done[node.Id] {
		return nil
	}
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
		cycles = append(cycles, g.removeCycles(toNode, callstack, done)...)
	}
	done[node.Id] = true
	callstack.Pop()
	return cycles
}

func (g *Graph[T]) RemoveCycles(node *Node[T]) []Cycle {
	return g.removeCycles(node, utils.NewCallStack(), map[string]bool{})
}

func (g *Graph[T]) RemoveJohnsonCycles() []Cycle {
	johnsonCycles := topo.DirectedCyclesIn(g)
	cycles := make([]Cycle, len(johnsonCycles))
	for i, c := range johnsonCycles {
		stack := make([]string, len(c))
		for i, n := range c {
			stack[i] = n.(*Node[T]).Id
		}
		g.RemoveFromToEdge(stack[0], stack[1])
		cycles[i] = Cycle{
			Cause: [2]string{stack[0], stack[1]},
			Stack: stack,
		}
	}

	return cycles
}
