package graph

import (
	"hash/fnv"

	"gonum.org/v1/gonum/graph"

	"github.com/gabotechs/dep-tree/internal/utils"
)

type Node[T any] struct {
	// Id identifies the node with a string. This is typically the absolute path
	//  of a file in dep-tree.
	Id string
	// Errors This node might hold some errors that are worth rendering to the user.
	//  For example, if the node is a file, maybe it failed to be parsed.
	Errors []error
	// Data is a generic implementation-defined data bucket. Implementations can put
	//  whatever they want here.
	Data T
}

func hash(s string) int64 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return int64(h.Sum32())
}

var hashCached = utils.Cached(hash)

func MakeNode[T any](id string, data T) *Node[T] {
	return &Node[T]{
		Id:     id,
		Errors: make([]error, 0),
		Data:   data,
	}
}

func (n *Node[T]) AddErrors(err ...error) {
	n.Errors = append(n.Errors, err...)
}

func (n *Node[T]) ID() int64 {
	return hashCached(n.Id)
}

type Nodes[T any] struct {
	nodes []*Node[T]
	cur   int
}

func NewNodesIterator[T any](nodes []*Node[T]) *Nodes[T] {
	return &Nodes[T]{
		nodes: nodes,
		cur:   -1,
	}
}

func (n *Nodes[T]) Next() bool {
	n.cur += 1
	return n.cur < len(n.nodes)
}

func (n *Nodes[T]) Len() int {
	return len(n.nodes)
}

func (n *Nodes[T]) Reset() {
	n.cur = -1
}

func (n *Nodes[T]) Node() graph.Node {
	if n.cur < 0 || n.cur >= len(n.nodes) {
		return nil
	} else {
		return n.nodes[n.cur]
	}
}
