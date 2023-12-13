package dep_tree

import (
	"context"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/utils"
)

const (
	testDir = ".render_test"
)

func TestRenderGraph(t *testing.T) {
	tests := []struct {
		Name string
		Spec [][]int
	}{
		{
			Name: "Simple",
			Spec: [][]int{
				{1, 2, 3},
				{2, 4},
				{3, 4},
				{4},
				{3},
			},
		},
		{
			Name: "Two in the same level",
			Spec: [][]int{
				{1, 2, 3},
				{3},
				{3},
				{},
			},
		},
		{
			Name: "Cyclic deps",
			Spec: [][]int{
				{1},
				{2},
				{1},
			},
		},
		{
			Name: "Children and Parents should be consistent",
			Spec: [][]int{
				{1, 2},
				{},
				{1},
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
				{1, 2, 3},
				{2, 4, 4275},
				{3, 4},
				{1423},
				{},
			},
		},
	}

	for _, tt := range tests {
		t.Run("Graph render: "+tt.Name, func(t *testing.T) {
			a := require.New(t)
			testParser := TestParser{
				Start: "0",
				Spec:  tt.Spec,
			}

			_, dt, err := NewDepTree[[]int](context.Background(), &testParser)
			a.NoError(err)

			board, err := dt.Render(testParser.Display)
			a.NoError(err)
			result, err := board.Render()
			a.NoError(err)

			outFile := path.Join(testDir, path.Base(tt.Name+".txt"))
			utils.GoldenTest(t, outFile, result)
		})

		t.Run("Structured render"+tt.Name, func(t *testing.T) {
			a := require.New(t)

			rendered, err := PrintStructured(
				context.Background(),
				"0",
				func(ctx context.Context, s string) (context.Context, NodeParser[[]int], error) {
					return ctx, &TestParser{Start: s, Spec: tt.Spec}, nil
				},
			)
			a.NoError(err)

			renderOutFile := path.Join(testDir, path.Base(tt.Name+".json"))
			utils.GoldenTest(t, renderOutFile, rendered)
		})
	}
}
