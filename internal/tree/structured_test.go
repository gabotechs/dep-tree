package tree

import (
	"path/filepath"
	"testing"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/utils"
)

const (
	structuredDir = ".structured_test"
)

func TestDepTree_RenderStructuredGraph(t *testing.T) {
	tests := []struct {
		Name string
		Spec [][]int
	}{
		{
			Name: "Simple",
			Spec: [][]int{
				0: {1, 2, 3},
				1: {2, 4},
				2: {3, 4},
				3: {4},
				4: {3},
			},
		},
		{
			Name: "Two in the same level",
			Spec: [][]int{
				0: {1, 2, 3},
				1: {3},
				2: {3},
				3: {},
			},
		},
		{
			Name: "Cyclic deps",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {1},
			},
		},
		{
			Name: "Children and Parents should be consistent",
			Spec: [][]int{
				0: {1, 2},
				1: {},
				2: {1},
			},
		},
		{
			Name: "Weird cycle combination",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {3},
				3: {4},
				4: {2},
			},
		},
		{
			Name: "Some nodes have errors",
			Spec: [][]int{
				0: {1, 2, 3},
				1: {2, 4, 4275},
				2: {3, 4},
				3: {1423},
				4: {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			tree, err := NewTree[[]int](
				[]string{"0"},
				&graph.TestParser{Spec: tt.Spec},
				func(node *graph.Node[[]int]) string { return node.Id },
				nil,
			)
			a.NoError(err)

			rendered, err := tree.RenderStructured()
			a.NoError(err)

			renderOutFile := filepath.Join(structuredDir, filepath.Base(tt.Name+".json"))
			utils.GoldenTest(t, renderOutFile, rendered)
		})
	}
}
