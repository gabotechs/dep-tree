package dep_tree

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const RebuildTestsEnv = "REBUILD_TESTS"

const testDir = ".render_test"

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
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			testParser := TestParser{
				Start: "0",
				Spec:  tt.Spec,
			}

			ctx := context.Background()

			ctx, dt, err := NewDepTree[[]int](ctx, &testParser)
			a.NoError(err)

			_, board, err := dt.Render(ctx, testParser.Display)
			a.NoError(err)
			result, err := board.Render()
			a.NoError(err)
			print(result)

			outFile := path.Join(testDir, path.Base(tt.Name+".txt"))
			if fileExists(outFile) && os.Getenv(RebuildTestsEnv) != "true" {
				expected, err := os.ReadFile(outFile)
				a.NoError(err)
				a.Equal(string(expected), result)
			} else {
				_ = os.Mkdir(testDir, os.ModePerm)
				err := os.WriteFile(outFile, []byte(result), os.ModePerm)
				a.NoError(err)
			}
		})
	}
}
