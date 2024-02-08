package graph

import (
	"gonum.org/v1/gonum/graph/topo"

	"github.com/gabotechs/dep-tree/internal/utils"
)

type Cycle struct {
	Cause [2]string
	Stack []string
}

func (g *Graph[T]) removeCyclesStartingFromNode(node *Node[T], callstack *utils.CallStack, done map[string]bool) []Cycle {
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
		cycles = append(cycles, g.removeCyclesStartingFromNode(toNode, callstack, done)...)
	}
	done[node.Id] = true
	callstack.Pop()
	return cycles
}

// RemoveCyclesStartingFromNode removes cycles that appear in order while performing a depth
// first search from a certain node. Note that this may not remove all cycles in the graph.
func (g *Graph[T]) RemoveCyclesStartingFromNode(node *Node[T]) []Cycle {
	return g.removeCyclesStartingFromNode(node, utils.NewCallStack(), map[string]bool{})
}

// RemoveElementaryCycles removes all the elementary cycles in the graph. The result
// of this can be non-deterministic.
func (g *Graph[T]) RemoveElementaryCycles() []Cycle {
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

// RemoveCycles removes all cycles in the graph, giving preference to cycles that start
// from the provided nodes.
func (g *Graph[T]) RemoveCycles(nodes []*Node[T]) []Cycle {
	var cycles []Cycle

	// First, remove the cycles computed from each entrypoint. This allows
	// us trim the cycles in a more "controlled way"
	for _, node := range nodes {
		for _, cycle := range g.RemoveCyclesStartingFromNode(node) {
			cycles = append(cycles, cycle)
		}
	}
	// Then, remove the cycles computed without taking entrypoints into account.
	// These are not as nice, as the rule for determining which cycles are trimmed is more arbitrary.
	for _, cycle := range g.RemoveElementaryCycles() {
		cycles = append(cycles, cycle)
	}
	return cycles
}
