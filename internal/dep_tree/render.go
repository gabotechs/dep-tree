package dep_tree

import (
	"context"
	"strconv"

	"dep-tree/internal/board"
	"dep-tree/internal/graph"
	"dep-tree/internal/utils"
)

const indent = 2
const NodeIdTag = "nodeId"
const NodeIndexTag = "nodeIndex"
const ConnectorOriginNodeIdTag = "connectorOrigin"
const ConnectorDestinationNodeIdTag = "connectorDestination"
const NodeParentsTag = "nodeParents"

func (dt *DepTree[T]) Render(
	ctx context.Context,
	display func(node *graph.Node[T]) string,
) (context.Context, *board.Board, error) {
	b := board.MakeBoard()

	lastLevel := -1
	prefix := ""
	xOffsetCount := 0
	xOffset := 0
	yOffset := 0
	for i, n := range dt.Nodes {
		if n.Lvl == lastLevel {
			if len(dt.Graph.Children(dt.Nodes[i-1].Node.Id)) > 0 {
				xOffsetCount++
				prefix += " "
			}
		} else {
			lastLevel = n.Lvl
			xOffset += xOffsetCount
			xOffsetCount = 0
			prefix = ""
			if i != 0 {
				yOffset++
			}
		}

		parents := dt.Graph.Parents(n.Node.Id)

		tags := map[string]string{
			NodeIdTag:      n.Node.Id,
			NodeIndexTag:   strconv.Itoa(i),
			NodeParentsTag: "",
		}

		for _, parent := range parents {
			tags[NodeParentsTag] += parent.Id + ";"
		}

		err := b.AddBlock(
			&board.Block{
				Id:       n.Node.Id,
				Label:    prefix + display(n.Node),
				Position: utils.Vec(indent*n.Lvl+xOffset, i+yOffset),
				Tags:     tags,
			},
		)
		if err != nil {
			return ctx, nil, err
		}
	}

	for _, n := range dt.Nodes {
		for _, child := range dt.Graph.Children(n.Node.Id) {
			tags := map[string]string{
				ConnectorOriginNodeIdTag:      n.Node.Id,
				ConnectorDestinationNodeIdTag: child.Id,
			}

			err := b.AddConnector(n.Node.Id, child.Id, tags)
			if err != nil {
				return ctx, nil, err
			}
		}
	}
	return ctx, b, nil
}
