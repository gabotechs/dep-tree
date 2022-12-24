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
	fr *node.FamilyRegistry[T],
) (context.Context, error) {
	if fr.Has(root.Id) {
		return ctx, nil
	}

	ctx, deps, err := parser.Deps(ctx, root)
	if err != nil {
		return ctx, err
	}

	fr.AddNode(root)

	for _, dep := range deps {
		ctx, err = computeDeps(ctx, dep, parser, fr)
		if err != nil {
			return ctx, err
		}
		err = fr.AddChild(root.Id, dep.Id)
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

func sortNodes[T any](
	ctx context.Context,
	root *node.Node[T],
	fr *node.FamilyRegistry[T],
) (context.Context, []graphNode[T]) {
	allNodes := fr.Nodes()
	result := make([]graphNode[T], 0)
	for _, n := range allNodes {
		var level int
		ctx, level = fr.Level(ctx, n.Id, root.Id)
		result = append(result, graphNode[T]{node: n, level: level})
	}

	sort.SliceStable(result, func(i, j int) bool {
		if result[i].level == result[j].level {
			return result[i].node.Id < result[j].node.Id
		} else {
			return result[i].level < result[j].level
		}
	})
	return ctx, result
}

const indent = 2

func renderGraph[T any](
	ctx context.Context,
	root *node.Node[T],
	parser NodeParser[T],
	fr *node.FamilyRegistry[T],
	nodes []graphNode[T],
) (context.Context, string, error) {
	b := board.MakeBoard()

	lastLevel := -1
	prefix := ""
	xOffsetCount := 0
	xOffset := 0
	for i, n := range nodes {
		if n.level == lastLevel {
			if len(fr.Children(nodes[i-1].node.Id)) > 0 {
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
		for _, child := range fr.Children(n.node.Id) {
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
	fr := node.NewFamilyRegistry[T]()
	ctx, err = computeDeps(ctx, root, parser, fr)
	if err != nil {
		return ctx, "", err
	}
	ctx, nodes := sortNodes(ctx, root, fr)
	return renderGraph(ctx, root, parser, fr, nodes)
}
