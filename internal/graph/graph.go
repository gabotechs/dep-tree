package graph

import (
	"context"
	"dep-tree/internal/node"
	"dep-tree/internal/render"
	"sort"
)

type NodeParser[T any] interface {
	Display(node *node.Node[T]) string
	Parse(id string) (*node.Node[T], error)
	Deps(node *node.Node[T]) []string
}

func makeGraph[T any](
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
		child, err := makeGraph(dep, parser, seen)
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

func renderGraph[T any](
	root *node.Node[T],
	parser NodeParser[T],
	rootId string,
) (string, error) {
	ctx := context.Background()
	allNodes := root.Flatten()

	nodesWithLevel := make([]graphNode[T], 0)
	maxLevel := 0
	maxSize := 0
	for _, k := range allNodes.Keys() {
		n, _ := allNodes.Get(k)
		var level int
		ctx, level = n.Level(ctx, rootId)
		if level > maxLevel {
			maxLevel = level
		}
		display := parser.Display(n)
		if len(display) > maxSize {
			maxSize = len(display)
		}
		nodesWithLevel = append(nodesWithLevel, graphNode[T]{node: n, level: level})
	}

	sort.SliceStable(nodesWithLevel, func(i, j int) bool {
		if nodesWithLevel[i].level == nodesWithLevel[j].level {
			return nodesWithLevel[i].node.Id < nodesWithLevel[j].node.Id
		} else {
			return nodesWithLevel[i].level < nodesWithLevel[j].level
		}
	})

	board := render.MakeBoard(render.BoardOptions{
		Indent:    2,
		BlockSize: maxSize,
	})

	for i, nodeWithLevel := range nodesWithLevel {
		err := board.AddBlock(
			nodeWithLevel.node.Id,
			parser.Display(nodeWithLevel.node),
			nodeWithLevel.level,
			i,
		)
		if err != nil {
			return "", err
		}
	}

	for _, nodeWithLevel := range nodesWithLevel {
		for _, childId := range nodeWithLevel.node.Children.Keys() {
			child, _ := nodeWithLevel.node.Children.Get(childId)
			err := board.AddDep(nodeWithLevel.node.Id, child.Id)
			if err != nil {
				return "", err
			}
		}
	}

	return board.Render(), nil
}

func RenderGraph[T any](
	entrypoint string,
	parser NodeParser[T],
) (string, error) {
	root, err := makeGraph(entrypoint, parser, map[string]*node.Node[T]{})
	if err != nil {
		return "", err
	}
	return renderGraph(root, parser, root.Id)
}
