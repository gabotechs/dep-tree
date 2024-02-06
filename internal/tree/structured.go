package tree

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/utils"
)

type StructuredTree struct {
	Tree                 map[string]interface{} `json:"tree" yaml:"tree"`
	CircularDependencies [][]string             `json:"circularDependencies" yaml:"circularDependencies"`
	Errors               map[string][]string    `json:"errors" yaml:"errors"`
}

func (t *Tree[T]) makeStructuredTree(
	from string,
	stack *utils.CallStack,
	cache map[string]map[string]interface{},
) (map[string]interface{}, error) {
	if stack == nil {
		stack = utils.NewCallStack()
	}
	if cache == nil {
		cache = make(map[string]map[string]interface{})
	}
	if result, ok := cache[from]; ok {
		return result, nil
	}
	err := stack.Push(from)
	if err != nil {
		return nil, errors.New("cannot create a structured object out of the graph because it contains at least 1 cycle: " + err.Error())
	}

	var result map[string]interface{}
	for _, to := range t.Graph.FromId(from) {
		if result == nil {
			result = make(map[string]interface{})
		}
		var err error
		result[t.NodeParser.Display(to).Name], err = t.makeStructuredTree(to.Id, stack, cache)
		if err != nil {
			return nil, err
		}
	}
	stack.Pop()
	cache[from] = result
	return result, nil
}

func (t *Tree[T]) RenderStructured() ([]byte, error) {
	if len(t.Entrypoints) > 1 {
		return nil, fmt.Errorf("this functionality requires that only 1 entrypoint is provided, but %d where detected. Consider providing a single entrypoint to your program", len(t.Entrypoints))
	}

	println("building structured tree")
	tree, err := t.makeStructuredTree(t.Entrypoints[0].Id, nil, nil)
	if err != nil {
		return nil, err
	}
	println("structured tree built")

	structuredTree := StructuredTree{
		Tree: map[string]interface{}{
			t.NodeParser.Display(t.Entrypoints[0]).Name: tree,
		},
		CircularDependencies: make([][]string, 0),
		Errors:               make(map[string][]string),
	}

	for _, cycle := range t.Cycles.Keys() {
		cycleDep, _ := t.Cycles.Get(cycle)
		renderedCycle := make([]string, len(cycleDep.Stack))
		for i, cycleDepEntry := range cycleDep.Stack {
			renderedCycle[i] = t.NodeParser.Display(t.Graph.Get(cycleDepEntry)).Name
		}
		structuredTree.CircularDependencies = append(structuredTree.CircularDependencies, renderedCycle)
	}

	for _, node := range t.Nodes {
		if node.Node.Errors != nil && len(node.Node.Errors) > 0 {
			erroredNode := t.NodeParser.Display(t.Graph.Get(node.Node.Id)).Name
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
	parser dep_tree.NodeParser[T],
) (string, error) {
	dt := dep_tree.NewDepTree(parser, files)
	err := dt.LoadGraph()
	if err != nil {
		return "", err
	}
	dt.LoadCycles()

	tree, err := NewTree(dt)
	if err != nil {
		return "", err
	}

	output, err := tree.RenderStructured()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
