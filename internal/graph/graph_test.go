package graph

import (
	"dep-tree/internal/node"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const RebuildTestsEnv = "REBUILD_TESTS"

type TestGraph struct {
	Spec [][]int
}

var _ NodeParser[[]int] = &TestGraph{}

func (t *TestGraph) Parse(id string) (*node.Node[[]int], error) {
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
	return node.MakeNode(id, children), nil
}

func (t *TestGraph) Deps(n *node.Node[[]int]) []string {
	result := make([]string, 0)
	for _, child := range n.Data {
		result = append(result, strconv.Itoa(child))
	}
	return result
}

func (t *TestGraph) Display(n *node.Node[[]int]) string {
	return n.Id
}

const testDir = "graph_test"

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func TestMakeGraph(t *testing.T) {
	a := require.New(t)
	files, err := os.ReadDir(testDir)
	a.NoError(err)
	for _, file := range files {
		if file.IsDir() || !strings.Contains(file.Name(), ".json") {
			continue
		}

		t.Run(file.Name(), func(t *testing.T) {
			a := require.New(t)
			content, err := os.ReadFile(path.Join(testDir, file.Name()))
			a.NoError(err)
			var testParser TestGraph
			err = json.Unmarshal(content, &testParser.Spec)
			a.NoError(err)

			result, err := RenderGraph[[]int]("0", &testParser)
			a.NoError(err)

			outFile := path.Join(testDir, strings.ReplaceAll(file.Name(), ".json", ".txt"))
			if fileExists(outFile) && os.Getenv(RebuildTestsEnv) != "true" {
				expected, err := os.ReadFile(outFile)
				a.NoError(err)
				a.Equal(string(expected), result)
			} else {
				err := os.WriteFile(outFile, []byte(result), os.ModePerm)
				a.NoError(err)
			}
		})
	}
}
