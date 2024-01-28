package graph

import (
	"fmt"

	om "github.com/elliotchance/orderedmap/v2"
	"gonum.org/v1/gonum/graph"
)

type Graph[T any] struct {
	nodes *om.OrderedMap[int64, *Node[T]]
	// Here "from" means: from node X I can reach nodes A, B and C
	// file -> dep
	fromEdges *om.OrderedMap[int64, *om.OrderedMap[int64, bool]]
	// Here "to" means: node X can be reached by A, B and C
	// dep -> file
	toEdges *om.OrderedMap[int64, *om.OrderedMap[int64, bool]]
}

var _ graph.Directed = &Graph[any]{}

func (g *Graph[T]) Node(id int64) graph.Node {
	v, _ := g.nodes.Get(id)
	return v
}

func (g *Graph[T]) Nodes() graph.Nodes {
	return NewNodesIterator(g.AllNodes())
}

func (g *Graph[T]) From(id int64) graph.Nodes {
	return NewNodesIterator(g.from(id))
}

func (g *Graph[T]) HasEdgeBetween(xid, yid int64) bool {
	if toNodes, ok := g.fromEdges.Get(xid); ok {
		if _, ok = toNodes.Get(yid); ok {
			return true
		}
	}
	if fromNodes, ok := g.toEdges.Get(xid); ok {
		if _, ok = fromNodes.Get(yid); ok {
			return true
		}
	}
	return false
}

func (g *Graph[T]) Edge(uid, vid int64) graph.Edge {
	if toNodes, ok := g.fromEdges.Get(uid); ok {
		if _, ok := toNodes.Get(vid); ok {
			if uNode, ok := g.nodes.Get(uid); ok {
				if vNode, ok := g.nodes.Get(vid); ok {
					return &Edge[T]{from: uNode, to: vNode}
				} else {
					panic(fmt.Sprintf("there was an Edge from %d to %d, but to node did not exist", uid, vid))
				}
			} else {
				panic(fmt.Sprintf("there was an Edge from %d to %d, but from node did not exist", uid, vid))
			}
		}
	}

	return nil
}

func (g *Graph[T]) HasEdgeFromTo(uid, vid int64) bool {
	if toNodes, ok := g.fromEdges.Get(uid); ok {
		if _, ok := toNodes.Get(vid); ok {
			return true
		}
	}
	return false
}

func (g *Graph[T]) To(id int64) graph.Nodes {
	return NewNodesIterator(g.to(id))
}

func NewGraph[T any]() *Graph[T] {
	return &Graph[T]{
		nodes:     om.NewOrderedMap[int64, *Node[T]](),
		fromEdges: om.NewOrderedMap[int64, *om.OrderedMap[int64, bool]](),
		toEdges:   om.NewOrderedMap[int64, *om.OrderedMap[int64, bool]](),
	}
}

func (g *Graph[T]) Has(nodeId string) bool {
	_, ok := g.nodes.Get(hashCached(nodeId))
	return ok
}

func (g *Graph[T]) AddNode(node *Node[T]) {
	g.nodes.Set(node.ID(), node)
}

func (g *Graph[T]) AddFromToEdge(fromId string, toIds ...string) error {
	from := hashCached(fromId)
	if _, ok := g.nodes.Get(from); !ok {
		return fmt.Errorf("'%s' is not in graph", fromId)
	}

	for _, toId := range toIds {
		to := hashCached(toId)
		if _, ok := g.nodes.Get(to); !ok {
			return fmt.Errorf("'%s' is not in graph", toId)
		}
		if toNodes, ok := g.fromEdges.Get(from); ok {
			toNodes.Set(to, true)
		} else {
			toNodes = om.NewOrderedMap[int64, bool]()
			toNodes.Set(to, true)
			g.fromEdges.Set(from, toNodes)
		}
		if fromNodes, ok := g.toEdges.Get(to); ok {
			fromNodes.Set(from, true)
		} else {
			fromNodes = om.NewOrderedMap[int64, bool]()
			fromNodes.Set(from, true)
			g.toEdges.Set(to, fromNodes)
		}
	}
	return nil
}

func (g *Graph[T]) RemoveFromToEdge(fromId string, toId string) {
	from := hashCached(fromId)
	to := hashCached(toId)
	if toNodes, ok := g.fromEdges.Get(from); ok {
		toNodes.Delete(to)
	}
	if fromNodes, ok := g.toEdges.Get(to); ok {
		fromNodes.Delete(from)
	}
}

func (g *Graph[T]) AllNodes() []*Node[T] {
	result := make([]*Node[T], g.nodes.Len())
	for i, nodeId := range g.nodes.Keys() {
		node, _ := g.nodes.Get(nodeId)
		result[i] = node
	}
	return result
}

func (g *Graph[T]) Get(id string) *Node[T] {
	node, _ := g.nodes.Get(hashCached(id))
	return node
}

// FromId returns the nodes to which id can reach.
func (g *Graph[T]) FromId(id string) []*Node[T] {
	return g.from(hashCached(id))
}

func (g *Graph[T]) from(idHash int64) []*Node[T] {
	if toNodes, ok := g.fromEdges.Get(idHash); ok {
		result := make([]*Node[T], toNodes.Len())
		for i, to := range toNodes.Keys() {
			if toNode, ok := g.nodes.Get(to); ok {
				result[i] = toNode
			}
		}
		return result
	} else {
		return make([]*Node[T], 0)
	}
}

// ToId returns the nodes from which id is reachable.
func (g *Graph[T]) ToId(id string) []*Node[T] {
	return g.to(hashCached(id))
}

func (g *Graph[T]) to(idHash int64) []*Node[T] {
	if fromNodes, ok := g.toEdges.Get(idHash); ok {
		result := make([]*Node[T], fromNodes.Len())
		for i, from := range fromNodes.Keys() {
			if fromNode, ok := g.nodes.Get(from); ok {
				result[i] = fromNode
			}
		}
		return result
	} else {
		return make([]*Node[T], 0)
	}
}

func (g *Graph[T]) GetNodesWithoutParents() []*Node[T] {
	result := make([]*Node[T], 0)
	for el := g.nodes.Front(); el != nil; el = el.Next() {
		if nodes, ok := g.toEdges.Get(el.Key); ok {
			if nodes.Len() == 0 {
				result = append(result, el.Value)
			}
		} else {
			result = append(result, el.Value)
		}
	}
	return result
}
