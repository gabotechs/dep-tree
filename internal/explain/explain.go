package explain

import (
	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/utils"
)

func Explain[T any](
	parser graph.NodeParser[T],
	fromFiles []string,
	toFiles []string,
	display func(node *graph.Node[T]) string,
	callbacks graph.LoadCallbacks[T],
) ([]string, error) {
	// 1. Build the graph.
	g := graph.NewGraph[T]()
	err := g.Load(append(fromFiles, toFiles...), parser, callbacks)
	if err != nil {
		return nil, err
	}

	// 2. Load all the dependencies between the two batches of files.
	fromSet := utils.SetFromSlice(fromFiles)

	nodes := g.AllNodes()
	var deps [][2]*graph.Node[T]
	for _, node := range nodes {
		if fromSet.Has(node.Id) {
			for _, toFile := range toFiles {
				toNode := g.Get(toFile)
				if toNode != nil && g.HasEdgeFromTo(node.ID(), toNode.ID()) {
					deps = append(deps, [2]*graph.Node[T]{node, toNode})
				}
			}
		}
	}

	// 3. Display the dependencies.
	result := make([]string, len(deps))
	for i, dep := range deps {
		result[i] = display(dep[0]) + " -> " + display(dep[1])
	}

	return result, nil
}
