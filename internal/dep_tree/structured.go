package dep_tree

import (
	"encoding/json"
	"errors"
	"fmt"

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
		result[dt.NodeParser.Display(to).Name], err = dt.makeStructuredTree(to.Id, nil)
		if err != nil {
			return nil, err
		}
	}
	stack.Pop()
	return result, nil
}

func (dt *DepTree[T]) RenderStructured() ([]byte, error) {
	if len(dt.Entrypoints) > 1 {
		return nil, fmt.Errorf("this functionality requires that only 1 entrypoint is provided, but %d where detected. Consider providing a single entrypoint to your program", len(dt.Entrypoints))
	}

	tree, err := dt.makeStructuredTree(dt.Entrypoints[0].Id, nil)
	if err != nil {
		return nil, err
	}

	structuredTree := StructuredTree{
		Tree: map[string]interface{}{
			dt.NodeParser.Display(dt.Entrypoints[0]).Name: tree,
		},
		CircularDependencies: make([][]string, 0),
		Errors:               make(map[string][]string),
	}

	for _, cycle := range dt.Cycles.Keys() {
		cycleDep, _ := dt.Cycles.Get(cycle)
		renderedCycle := make([]string, len(cycleDep.Stack))
		for i, cycleDepEntry := range cycleDep.Stack {
			renderedCycle[i] = dt.NodeParser.Display(dt.Graph.Get(cycleDepEntry)).Name
		}
		structuredTree.CircularDependencies = append(structuredTree.CircularDependencies, renderedCycle)
	}

	for _, node := range dt.Nodes {
		if node.Node.Errors != nil && len(node.Node.Errors) > 0 {
			erroredNode := dt.NodeParser.Display(dt.Graph.Get(node.Node.Id)).Name
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
	files []string,
	parser NodeParser[T],
) (string, error) {
	dt := NewDepTree(parser, files)
	err := dt.LoadDeps()
	if err != nil {
		return "", err
	}
	output, err := dt.RenderStructured()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
