package node

import (
	"fmt"

	om "github.com/elliotchance/orderedmap/v2"
)

type Node[T any] struct {
	Id   string
	Data T
}

func MakeNode[T any](id string, data T) *Node[T] {
	return &Node[T]{
		Id:   id,
		Data: data,
	}
}

type FamilyRegistry[T any] struct {
	nodes       *om.OrderedMap[string, *Node[T]]
	childEdges  *om.OrderedMap[string, *om.OrderedMap[string, bool]]
	parentEdges *om.OrderedMap[string, *om.OrderedMap[string, bool]]
}

func NewFamilyRegistry[T any]() *FamilyRegistry[T] {
	return &FamilyRegistry[T]{
		nodes:       om.NewOrderedMap[string, *Node[T]](),
		childEdges:  om.NewOrderedMap[string, *om.OrderedMap[string, bool]](),
		parentEdges: om.NewOrderedMap[string, *om.OrderedMap[string, bool]](),
	}
}

func (f *FamilyRegistry[T]) Has(nodeId string) bool {
	_, ok := f.nodes.Get(nodeId)
	return ok
}

func (f *FamilyRegistry[T]) AddNode(node *Node[T]) {
	f.nodes.Set(node.Id, node)
}

func (f *FamilyRegistry[T]) AddChild(parent string, children ...string) error {
	if _, ok := f.nodes.Get(parent); !ok {
		return fmt.Errorf("'%s' is not in graph", parent)
	}
	for _, child := range children {
		if _, ok := f.nodes.Get(child); !ok {
			return fmt.Errorf("'%s' is not in graph", child)
		}
	}

	for _, child := range children {
		if childEdges, ok := f.childEdges.Get(parent); ok {
			childEdges.Set(child, true)
		} else {
			childEdges := om.NewOrderedMap[string, bool]()
			childEdges.Set(child, true)
			f.childEdges.Set(parent, childEdges)
		}
		if parentEdges, ok := f.parentEdges.Get(child); ok {
			parentEdges.Set(parent, true)
		} else {
			parentEdges := om.NewOrderedMap[string, bool]()
			parentEdges.Set(parent, true)
			f.parentEdges.Set(child, parentEdges)
		}
	}
	return nil
}

func (f *FamilyRegistry[T]) Nodes() []*Node[T] {
	result := make([]*Node[T], f.nodes.Len())
	for i, nodeId := range f.nodes.Keys() {
		node, _ := f.nodes.Get(nodeId)
		result[i] = node
	}
	return result
}

func (f *FamilyRegistry[T]) Children(id string) []*Node[T] {
	if children, ok := f.childEdges.Get(id); ok {
		result := make([]*Node[T], children.Len())
		for i, child := range children.Keys() {
			if childNode, ok := f.nodes.Get(child); ok {
				result[i] = childNode
			}
		}
		return result
	} else {
		return make([]*Node[T], 0)
	}
}

func (f *FamilyRegistry[T]) Parents(id string) []*Node[T] {
	if parents, ok := f.parentEdges.Get(id); ok {
		result := make([]*Node[T], parents.Len())
		for i, parent := range parents.Keys() {
			if parentNode, ok := f.nodes.Get(parent); ok {
				result[i] = parentNode
			}
		}
		return result
	} else {
		return make([]*Node[T], 0)
	}
}
