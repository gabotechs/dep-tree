package graph

type Node[T any] struct {
	Id     string
	Errors []error
	Data   T
}

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
