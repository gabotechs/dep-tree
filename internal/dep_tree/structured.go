package dep_tree

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gabotechs/dep-tree/internal/utils"
)

type StructuredTree struct {
	Tree                 map[string]interface{} `json:"tree" yaml:"tree"`
	CircularDependencies [][]string             `json:"circularDependencies" yaml:"circularDependencies"`
	Errors               map[string][]string    `json:"errors" yaml:"errors"`
}

func (dt *DepTree[T]) makeStructuredTree(
	from string,
	stack *utils.CallStack,
) (map[string]interface{}, error) {
	if stack == nil {
		stack = utils.NewCallStack()
	}
	err := stack.Push(from)
	if err != nil {
		return nil, errors.New("cannot create a structured object out of the graph because it contains at least 1 cycle: " + err.Error())
	}
	var result map[string]interface{}
	for _, to := range dt.Graph.FromId(from) {
		if result == nil {
			result = make(map[string]interface{})
		}
		var err error
		result[dt.NodeParser.Display(to)], err = dt.makeStructuredTree(to.Id, nil)
		if err != nil {
			return nil, err
		}
	}
	stack.Pop()
	return result, nil
}

func (dt *DepTree[T]) RenderStructured() ([]byte, error) {
	root, err := dt.Root()
	if err != nil {
		return nil, err
	}

	tree, err := dt.makeStructuredTree(root.Id, nil)
	if err != nil {
		return nil, err
	}

	structuredTree := StructuredTree{
		Tree: map[string]interface{}{
			dt.NodeParser.Display(root): tree,
		},
		CircularDependencies: make([][]string, 0),
		Errors:               make(map[string][]string),
	}

	for _, cycle := range dt.Cycles.Keys() {
		cycleDep, _ := dt.Cycles.Get(cycle)
		renderedCycle := make([]string, len(cycleDep.Stack))
		for i, cycleDepEntry := range cycleDep.Stack {
			renderedCycle[i] = dt.NodeParser.Display(dt.Graph.Get(cycleDepEntry))
		}
		structuredTree.CircularDependencies = append(structuredTree.CircularDependencies, renderedCycle)
	}

	for _, node := range dt.Nodes {
		if node.Node.Errors != nil && len(node.Node.Errors) > 0 {
			erroredNode := dt.NodeParser.Display(dt.Graph.Get(node.Node.Id))
			nodeErrors := make([]string, len(node.Node.Errors))
			for i, err := range node.Node.Errors {
				nodeErrors[i] = err.Error()
			}
			structuredTree.Errors[erroredNode] = nodeErrors
		}
	}

	return json.MarshalIndent(structuredTree, "", "  ")
}

func PrintStructured[T any](
	ctx context.Context,
	entrypoint string,
	parserBuilder NodeParserBuilder[T],
) (string, error) {
	ctx, parser, err := parserBuilder(ctx, entrypoint)
	if err != nil {
		return "", err
	}
	dt := NewDepTree(parser)
	if err != nil {
		return "", err
	}
	_, err = dt.LoadDeps(ctx)
	if err != nil {
		return "", err
	}
	output, err := dt.RenderStructured()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
