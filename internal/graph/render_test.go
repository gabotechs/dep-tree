package graph

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

const RebuildTestsEnv = "REBUILD_TESTS"

type TestGraph struct {
	Spec [][]int
}

var _ NodeParser[[]int] = &TestGraph{}

func (t *TestGraph) Entrypoint(id string) (*Node[[]int], error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	var children []int
	if idInt >= len(t.Spec) {
		return nil, fmt.Errorf("%s not present in spec", id)
	} else {
		children = t.Spec[idInt]
	}
	return MakeNode(id, children), nil
}

func (t *TestGraph) Deps(ctx context.Context, n *Node[[]int]) (context.Context, []*Node[[]int], error) {
	result := make([]*Node[[]int], 0)
	for _, child := range n.Data {
		c, _ := t.Entrypoint(strconv.Itoa(child))
		result = append(result, c)
	}
	return ctx, result, nil
}

func (t *TestGraph) Display(n *Node[[]int]) string {
	return n.Id
}

const testDir = ".graph_test"

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
			testParser := TestGraph{
				Spec: tt.Spec,
			}

			ctx := context.Background()

			_, result, err := RenderGraph[[]int](ctx, "0", &testParser)
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
