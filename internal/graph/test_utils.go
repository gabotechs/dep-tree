package graph

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gammazero/deque"
)

func MakeTestGraph(spec [][]int) *Graph[int] {
	g := NewGraph[int]()

	var queue deque.Deque[*Node[int]]
	node := MakeNode("0", 0)
	g.AddNode(node)
	queue.PushBack(node)
	visited := make(map[string]bool)

	for queue.Len() > 0 {
		node := queue.PopFront()
		if _, ok := visited[node.Id]; ok {
			continue
		}
		visited[node.Id] = true
		deps := spec[node.Data]

		for _, dep := range deps {
			depId := strconv.Itoa(dep)
			depNode := g.Get(depId)
			if depNode == nil {
				depNode = MakeNode(strconv.Itoa(dep), dep)
				g.AddNode(depNode)
			} else {

			}
			err := g.AddFromToEdge(node.Id, depNode.Id)
			if err != nil {
				panic(err)
			}
			queue.PushBack(depNode)
		}
	}
	return g
}

type TestParser struct {
	Spec [][]int
}

var _ NodeParser[[]int] = &TestParser{}

func (t *TestParser) Node(id string) (*Node[[]int], error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	if idInt >= len(t.Spec) {
		return nil, fmt.Errorf("%d not present in spec", idInt)
	} else {
		return MakeNode(id, t.Spec[idInt]), nil
	}
}

func (t *TestParser) Deps(n *Node[[]int]) ([]*Node[[]int], error) {
	result := make([]*Node[[]int], 0)
	for _, child := range n.Data {
		if child < 0 {
			return nil, errors.New("no negative children")
		}
		c, err := t.Node(strconv.Itoa(child))
		if err != nil {
			n.Errors = append(n.Errors, err)
		} else {
			result = append(result, c)
		}
	}
	return result, nil
}
