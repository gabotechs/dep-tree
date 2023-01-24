package dep_tree

import (
	"context"

	"dep-tree/internal/graph"
)

func loadDeps[T any](
	ctx context.Context,
	g *graph.Graph[T],
	root *graph.Node[T],
	parser NodeParser[T],
) (context.Context, error) {
	if g.Has(root.Id) {
		return ctx, nil
	}

	ctx, deps, err := parser.Deps(ctx, root)
	if err != nil {
		return ctx, err
	}

	g.AddNode(root)

	for _, dep := range deps {
		ctx, err = loadDeps(ctx, g, dep, parser)
		if err != nil {
			return ctx, err
		}
		err = g.AddChild(root.Id, dep.Id)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}

func LoadDeps[T any](
	ctx context.Context,
	g *graph.Graph[T],
	parser NodeParser[T],
) (context.Context, string, error) {
	root, err := parser.Entrypoint()
	if err != nil {
		return ctx, "", err
	}
	ctx, err = loadDeps(ctx, g, root, parser)
	return ctx, root.Id, err
}
