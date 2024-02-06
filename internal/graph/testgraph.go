package graph

import (
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
