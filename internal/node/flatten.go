package node

import "github.com/elliotchance/orderedmap/v2"

func flatten[T any](root *Node[T], storage *orderedmap.OrderedMap[string, *Node[T]]) {
	if _, ok := storage.Get(root.Id); ok {
		return
	}
	storage.Set(root.Id, root)
	for _, childId := range root.Children.Keys() {
		child, _ := root.Children.Get(childId)
		flatten(child, storage)
	}
}

// Flatten retrieves a hashmap of unique non-repeated nodes
func (n *Node[T]) Flatten() *orderedmap.OrderedMap[string, *Node[T]] {
	storage := orderedmap.NewOrderedMap[string, *Node[T]]()
	flatten(n, storage)
	return storage
}
