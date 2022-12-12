package graph

import (
	"context"
	"sort"

	"dep-tree/internal/node"
	"dep-tree/internal/render"
)

type NodeParser[T any] interface {
	Display(node *node.Node[T]) string
	Parse(id string) (*node.Node[T], error)
	Deps(node *node.Node[T]) []string
}

func makeNodes[T any](
	entrypoint string,
	parser NodeParser[T],
	seen map[string]*node.Node[T],
) (*node.Node[T], error) {
	if cached, ok := seen[entrypoint]; ok {
		return cached, nil
	}
	root, err := parser.Parse(entrypoint)
	if err != nil {
		return nil, err
	} else if _, ok := seen[entrypoint]; !ok {
		seen[entrypoint] = root
	}

	deps := parser.Deps(root)
	for _, dep := range deps {
		child, err := makeNodes(dep, parser, seen)
		if err != nil {
			return nil, err
		}
		root.AddChild(child)
	}
	return root, nil
}

type graphNode[T any] struct {
	node  *node.Node[T]
	level int
}

func sortNodes[T any](root *node.Node[T]) []graphNode[T] {
	ctx := context.Background()
	allNodes := root.Flatten()

	result := make([]graphNode[T], 0)
	for _, k := range allNodes.Keys() {
		n, _ := allNodes.Get(k)
		var level int
		ctx, level = n.Level(ctx, root.Id)
		result = append(result, graphNode[T]{node: n, level: level})
	}

	sort.SliceStable(result, func(i, j int) bool {
		if result[i].level == result[j].level {
			return result[i].node.Id < result[j].node.Id
		} else {
			return result[i].level < result[j].level
		}
	})
	return result
}

const indent = 4

func renderGraph[T any](
	parser NodeParser[T],
	nodes []graphNode[T],
) (string, error) {
	board := render.MakeBoard()

	lastLevel := -1
	yOffset := 0
	for i, n := range nodes {
		if n.level == lastLevel {
			if nodes[i-1].node.Children.Len() > 0 {
				yOffset++
			}
		} else {
			lastLevel = n.level
		}

		err := board.AddBlock(
			n.node.Id,
			parser.Display(n.node),
			indent*n.level,
			i+yOffset,
		)
		if err != nil {
			return "", err
		}
	}

	for i := range nodes {
		n := nodes[len(nodes)-1-i]
		for _, childId := range n.node.Children.Keys() {
			child, _ := n.node.Children.Get(childId)
			err := board.AddConnector(n.node.Id, child.Id)
			if err != nil {
				return "", err
			}
		}
	}

	return board.Render()
}

func RenderGraph[T any](
	entrypoint string,
	parser NodeParser[T],
) (string, error) {
	root, err := makeNodes(entrypoint, parser, map[string]*node.Node[T]{})
	if err != nil {
		return "", err
	}
	nodes := sortNodes(root)
	return renderGraph(parser, nodes)
}
