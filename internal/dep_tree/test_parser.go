package dep_tree

import (
	"context"
	"fmt"
	"strconv"

	"dep-tree/internal/graph"
)

type TestParser struct {
	Start string
	Spec  [][]int
}

var _ NodeParser[[]int] = &TestParser{}

func (t *TestParser) getNode(id string) (*graph.Node[[]int], error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	var children []int
	if idInt >= len(t.Spec) {
		return nil, fmt.Errorf("%s not present in spec", t.Start)
	} else {
		children = t.Spec[idInt]
	}
	return graph.MakeNode(id, children), nil
}

func (t *TestParser) Entrypoint() (*graph.Node[[]int], error) {
	return t.getNode(t.Start)
}

func (t *TestParser) Deps(ctx context.Context, n *graph.Node[[]int]) (context.Context, []*graph.Node[[]int], error) {
	result := make([]*graph.Node[[]int], 0)
	for _, child := range n.Data {
		c, _ := t.getNode(strconv.Itoa(child))
		result = append(result, c)
	}
	return ctx, result, nil
}

func (t *TestParser) Display(n *graph.Node[[]int]) string {
	return n.Id
}
