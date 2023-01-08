package graph

import (
	"context"
	"sort"
	"strconv"

	"dep-tree/internal/board"
	"dep-tree/internal/utils"
)

type NodeParser[T any] interface {
	Display(node *Node[T]) string
	Entrypoint(entrypoint string) (*Node[T], error)
	Deps(ctx context.Context, node *Node[T]) (context.Context, []*Node[T], error)
}

func (g *Graph[T]) computeDeps(
	ctx context.Context,
	root *Node[T],
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
		ctx, err = g.computeDeps(ctx, dep, parser)
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

type graphNode[T any] struct {
	node  *Node[T]
	level int
}

func (g *Graph[T]) getSortNodes(
	ctx context.Context,
	root *Node[T],
) (context.Context, []graphNode[T]) {
	allNodes := g.Nodes()
	result := make([]graphNode[T], 0)
	for _, n := range allNodes {
		var level int
		ctx, level = g.Level(ctx, n.Id, root.Id)
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
const NodeIdTag = "nodeId"
const NodeIndexTag = "nodeIndex"
const ConnectorOriginNodeIdTag = "connectorOrigin"
const ConnectorDestinationNodeIdTag = "connectorDestination"
const NodeParentsTag = "nodeParents"

func (g *Graph[T]) renderGraph(
	ctx context.Context,
	parser NodeParser[T],
	nodes []graphNode[T],
) (context.Context, *board.Board, error) {
	b := board.MakeBoard()

	lastLevel := -1
	prefix := ""
	xOffsetCount := 0
	xOffset := 0
	yOffset := 0
	for i, n := range nodes {
		if n.level == lastLevel {
			if len(g.Children(nodes[i-1].node.Id)) > 0 {
				xOffsetCount++
				prefix += " "
			}
		} else {
			lastLevel = n.level
			xOffset += xOffsetCount
			xOffsetCount = 0
			prefix = ""
			if i != 0 {
				yOffset++
			}
		}

		parents := g.Parents(n.node.Id)

		tags := map[string]string{
			NodeIdTag:      n.node.Id,
			NodeIndexTag:   strconv.Itoa(i),
			NodeParentsTag: "",
		}

		for _, parent := range parents {
			tags[NodeParentsTag] += parent.Id + ";"
		}

		err := b.AddBlock(
			&board.Block{
				Id:       n.node.Id,
				Label:    prefix + parser.Display(n.node),
				Position: utils.Vec(indent*n.level+xOffset, i+yOffset),
				Tags:     tags,
			},
		)
		if err != nil {
			return ctx, nil, err
		}
	}

	for _, n := range nodes {
		for _, child := range g.Children(n.node.Id) {
			tags := map[string]string{
				ConnectorOriginNodeIdTag:      n.node.Id,
				ConnectorDestinationNodeIdTag: child.Id,
			}

			err := b.AddConnector(n.node.Id, child.Id, tags)
			if err != nil {
				return ctx, nil, err
			}
		}
	}
	return ctx, b, nil
}

func RenderGraph[T any](
	ctx context.Context,
	entrypoint string,
	parser NodeParser[T],
) (context.Context, *board.Board, error) {
	root, err := parser.Entrypoint(entrypoint)
	if err != nil {
		return ctx, nil, err
	}
	g := NewGraph[T]()
	ctx, err = g.computeDeps(ctx, root, parser)
	if err != nil {
		return ctx, nil, err
	}
	ctx, nodes := g.getSortNodes(ctx, root)
	return g.renderGraph(ctx, parser, nodes)
}
