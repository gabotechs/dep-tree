package dep_tree

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/utils"
)

const (
	renderDir = ".render_test"
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
			Name: "Weird cycle combination 2",
			Spec: [][]int{
				0: {3, 1},
				1: {2},
				2: {3},
				3: {4},
				4: {0},
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
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			testParser := TestParser{
				Start: "0",
				Spec:  tt.Spec,
			}

			dt := NewDepTree[[]int](&testParser)

			err := dt.LoadDeps()
			a.NoError(err)

			board, err := dt.Render()
			a.NoError(err)
			result, err := board.Render()
			a.NoError(err)

			outFile := filepath.Join(renderDir, path.Base(tt.Name+".txt"))
			utils.GoldenTest(t, outFile, result)
		})
	}
}
