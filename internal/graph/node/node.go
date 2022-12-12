package node

import "github.com/elliotchance/orderedmap/v2"

type Node[T any] struct {
	Id       string
	GroupId  string
	Data     T
	Children *orderedmap.OrderedMap[string, *Node[T]]
	Parents  *orderedmap.OrderedMap[string, *Node[T]]
}

func MakeNode[T any](id string, groupId string, data T) *Node[T] {
	return &Node[T]{
		Id:       id,
		Data:     data,
		GroupId:  groupId,
		Children: orderedmap.NewOrderedMap[string, *Node[T]](),
		Parents:  orderedmap.NewOrderedMap[string, *Node[T]](),
	}
}

func (n *Node[T]) AddChild(children ...*Node[T]) {
	for _, child := range children {
		n.Children.Set(child.Id, child)
		child.Parents.Set(n.Id, n)
	}
}
