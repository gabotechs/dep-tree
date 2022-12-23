package graph

import (
	"context"
	"sort"

	"dep-tree/internal/board"
	"dep-tree/internal/graph/node"
)

type NodeParser[T any] interface {
	Display(node *node.Node[T], root *node.Node[T]) string
	Entrypoint(entrypoint string) (*node.Node[T], error)
	Deps(ctx context.Context, node *node.Node[T]) (context.Context, []*node.Node[T], error)
}

func computeDeps[T any](
	ctx context.Context,
	root *node.Node[T],
	parser NodeParser[T],
	seen map[string]*node.Node[T],
) (context.Context, error) {
	if _, ok := seen[root.Id]; ok {
		return ctx, nil
	} else {
		seen[root.Id] = root
	}

	ctx, deps, err := parser.Deps(ctx, root)
	if err != nil {
		return ctx, err
	}
	for _, dep := range deps {
		root.AddChild(dep)
		ctx, err = computeDeps(ctx, dep, parser, seen)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
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

const indent = 2

func renderGraph[T any](
	ctx context.Context,
	parser NodeParser[T],
	root *node.Node[T],
	nodes []graphNode[T],
) (context.Context, string, error) {
	b := board.MakeBoard()

	lastLevel := -1
	prefix := ""
	xOffsetCount := 0
	xOffset := 0
	for i, n := range nodes {
		if n.level == lastLevel {
			if nodes[i-1].node.Children.Len() > 0 {
				xOffsetCount++
				prefix += " "
			}
		} else {
			lastLevel = n.level
			xOffset += xOffsetCount
			xOffsetCount = 0
			prefix = ""
		}

		err := b.AddBlock(
			n.node.Id,
			prefix+parser.Display(n.node, root),
			indent*n.level+xOffset,
			i,
		)
		if err != nil {
			return ctx, "", err
		}
	}

	for _, n := range nodes {
		for _, childId := range n.node.Children.Keys() {
			child, _ := n.node.Children.Get(childId)
			err := b.AddConnector(n.node.Id, child.Id)
			if err != nil {
				return ctx, "", err
			}
		}
	}
	rendered, err := b.Render()
	return ctx, rendered, err
}

func RenderGraph[T any](
	ctx context.Context,
	entrypoint string,
	parser NodeParser[T],
) (context.Context, string, error) {
	root, err := parser.Entrypoint(entrypoint)
	if err != nil {
		return ctx, "", err
	}
	ctx, err = computeDeps(ctx, root, parser, map[string]*node.Node[T]{})
	if err != nil {
		return ctx, "", err
	}
	nodes := sortNodes(root)
	return renderGraph(ctx, parser, root, nodes)
}
