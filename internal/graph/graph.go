package graph

import (
	"fmt"

	om "github.com/elliotchance/orderedmap/v2"
)

type Graph[T any] struct {
	nodes       *om.OrderedMap[string, *Node[T]]
	childEdges  *om.OrderedMap[string, *om.OrderedMap[string, bool]]
	parentEdges *om.OrderedMap[string, *om.OrderedMap[string, bool]]
}

func NewGraph[T any]() *Graph[T] {
	return &Graph[T]{
		nodes:       om.NewOrderedMap[string, *Node[T]](),
		childEdges:  om.NewOrderedMap[string, *om.OrderedMap[string, bool]](),
		parentEdges: om.NewOrderedMap[string, *om.OrderedMap[string, bool]](),
	}
}

func (g *Graph[T]) Has(nodeId string) bool {
	_, ok := g.nodes.Get(nodeId)
	return ok
}

func (g *Graph[T]) AddNode(node *Node[T]) {
	g.nodes.Set(node.Id, node)
}

func (g *Graph[T]) AddChild(parent string, children ...string) error {
	if _, ok := g.nodes.Get(parent); !ok {
		return fmt.Errorf("'%s' is not in graph", parent)
	}
	for _, child := range children {
		if _, ok := g.nodes.Get(child); !ok {
			return fmt.Errorf("'%s' is not in graph", child)
		}
	}

	for _, child := range children {
		if childEdges, ok := g.childEdges.Get(parent); ok {
			childEdges.Set(child, true)
		} else {
			childEdges = om.NewOrderedMap[string, bool]()
			childEdges.Set(child, true)
			g.childEdges.Set(parent, childEdges)
		}
		if parentEdges, ok := g.parentEdges.Get(child); ok {
			parentEdges.Set(parent, true)
		} else {
			parentEdges = om.NewOrderedMap[string, bool]()
			parentEdges.Set(parent, true)
			g.parentEdges.Set(child, parentEdges)
		}
	}
	return nil
}

func (g *Graph[T]) Nodes() []*Node[T] {
	result := make([]*Node[T], g.nodes.Len())
	for i, nodeId := range g.nodes.Keys() {
		node, _ := g.nodes.Get(nodeId)
		result[i] = node
	}
	return result
}

func (g *Graph[T]) Get(id string) *Node[T] {
	if node, ok := g.nodes.Get(id); ok {
		return node
	} else {
		return nil
	}
}

func (g *Graph[T]) Children(id string) []*Node[T] {
	if children, ok := g.childEdges.Get(id); ok {
		result := make([]*Node[T], children.Len())
		for i, child := range children.Keys() {
			if childNode, ok := g.nodes.Get(child); ok {
				result[i] = childNode
			}
		}
		return result
	} else {
		return make([]*Node[T], 0)
	}
}

func (g *Graph[T]) Parents(id string) []*Node[T] {
	if parents, ok := g.parentEdges.Get(id); ok {
		result := make([]*Node[T], parents.Len())
		for i, parent := range parents.Keys() {
			if parentNode, ok := g.nodes.Get(parent); ok {
				result[i] = parentNode
			}
		}
		return result
	} else {
		return make([]*Node[T], 0)
	}
}
