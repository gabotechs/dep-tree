package dep_tree

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const RebuildTestsEnv = "REBUILD_TESTS"

const (
	testDir            = ".render_test"
	RebuiltTestEnvTrue = "true"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

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
		t.Run(tt.Name, func(t *testing.T) {
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
			if fileExists(outFile) && os.Getenv(RebuildTestsEnv) != RebuiltTestEnvTrue {
				expected, err := os.ReadFile(outFile)
				a.NoError(err)
				a.Equal(string(expected), result)
			} else {
				_ = os.Mkdir(testDir, os.ModePerm)
				err := os.WriteFile(outFile, []byte(result), os.ModePerm)
				a.NoError(err)
			}

			rendered, err := dt.RenderStructured(testParser.Display)
			a.NoError(err)

			renderOutFile := path.Join(testDir, path.Base(tt.Name+".json"))
			if fileExists(renderOutFile) && os.Getenv(RebuildTestsEnv) != RebuiltTestEnvTrue {
				expected, err := os.ReadFile(renderOutFile)
				a.NoError(err)
				a.Equal(string(expected), string(rendered))
			} else {
				_ = os.Mkdir(testDir, os.ModePerm)
				err := os.WriteFile(renderOutFile, rendered, os.ModePerm)
				a.NoError(err)
			}
		})
	}
}
