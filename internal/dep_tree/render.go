package dep_tree

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gabotechs/dep-tree/internal/board"
	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/utils"
)

const indent = 2
const NodeIdTag = "nodeId"
const NodeIndexTag = "nodeIndex"
const ConnectorOriginNodeIdTag = "connectorOrigin"
const ConnectorDestinationNodeIdTag = "connectorDestination"
const NodeParentsTag = "nodeParents"

func (dt *DepTree[T]) Render(display func(node *graph.Node[T]) string) (*board.Board, error) {
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
			return nil, err
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
				return nil, err
			}
		}
	}
	return b, nil
}

type StructuredTree struct {
	Tree                 map[string]interface{} `json:"tree" yaml:"tree"`
	CircularDependencies [][]string             `json:"circularDependencies" yaml:"circularDependencies"`
	Errors               map[string][]string    `json:"errors" yaml:"errors"`
}

func (dt *DepTree[T]) makeStructuredTree(
	node string,
	display func(node *graph.Node[T]) string,
) map[string]interface{} {
	var result map[string]interface{}
	for _, child := range dt.Graph.Children(node) {
		if _, ok := dt.Cycles.Get([2]string{node, child.Id}); ok {
			continue
		}
		if result == nil {
			result = make(map[string]interface{})
		}
		result[display(child)] = dt.makeStructuredTree(child.Id, display)
	}
	return result
}

func (dt *DepTree[T]) RenderStructured(display func(node *graph.Node[T]) string) ([]byte, error) {
	root := dt.Graph.Get(dt.RootId)
	if root == nil {
		return nil, fmt.Errorf("could not retrieve root node from rootId %s", dt.RootId)
	}

	structuredTree := StructuredTree{
		Tree: map[string]interface{}{
			display(root): dt.makeStructuredTree(dt.RootId, display),
		},
		CircularDependencies: make([][]string, 0),
		Errors:               make(map[string][]string),
	}

	for _, cycle := range dt.Cycles.Keys() {
		cycleDep, _ := dt.Cycles.Get(cycle)
		renderedCycle := make([]string, len(cycleDep.Stack))
		for i, cycleDepEntry := range cycleDep.Stack {
			renderedCycle[i] = display(dt.Graph.Get(cycleDepEntry))
		}
		structuredTree.CircularDependencies = append(structuredTree.CircularDependencies, renderedCycle)
	}

	for _, node := range dt.Nodes {
		if node.Node.Errors != nil && len(node.Node.Errors) > 0 {
			erroredNode := display(dt.Graph.Get(node.Node.Id))
			errors := make([]string, len(node.Node.Errors))
			for i, err := range node.Node.Errors {
				errors[i] = err.Error()
			}
			structuredTree.Errors[erroredNode] = errors
		}
	}

	return json.MarshalIndent(structuredTree, "", "  ")
}

func PrintStructured[T any](
	ctx context.Context,
	entrypoint string,
	parserBuilder func(context.Context, string) (context.Context, NodeParser[T], error),
) (string, error) {
	ctx, parser, err := parserBuilder(ctx, entrypoint)
	if err != nil {
		return "", err
	}
	_, dt, err := NewDepTree(ctx, parser)
	if err != nil {
		return "", err
	}
	output, err := dt.RenderStructured(parser.Display)
	if err != nil {
		return "", err
	}
	return string(output), nil
}
