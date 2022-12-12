package board

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const RebuildTestsEnv = "REBUILD_TESTS"

const testPath = "board_test"

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func expectTest(t *testing.T, result string) {
	a := require.New(t)
	_ = os.Mkdir(testPath, os.ModePerm)
	fullPath := path.Join(testPath, path.Base(t.Name())+".txt")
	if fileExists(fullPath) && os.Getenv(RebuildTestsEnv) != "true" {
		expected, err := os.ReadFile(fullPath)
		a.NoError(err)
		a.Equal(string(expected), result)
	} else {
		print(result)
		err := os.WriteFile(fullPath, []byte(result), os.ModePerm)
		a.NoError(err)
	}
}

type TestBlock struct {
	name string
	x    int
	y    int
}
type TestConnection struct {
	from int
	to   []int
}

func TestBoard(t *testing.T) {
	tests := []struct {
		Name        string
		Blocks      []TestBlock
		Connections []TestConnection
	}{
		{
			Name: "SimpleDeps",
			Blocks: []TestBlock{
				{name: "index.ts", x: 0, y: 0},
				{name: "foo.ts", x: 3, y: 4},
				{name: "bar.ts", x: 5, y: 5},
			},
			Connections: []TestConnection{
				{from: 0, to: []int{1}},
				{from: 1, to: []int{2}},
			},
		},
		{
			Name: "ReverseDeps",
			Blocks: []TestBlock{
				{name: "a", x: 0, y: 0},
				{name: "b", x: 2, y: 1},
			},
			Connections: []TestConnection{
				{from: 1, to: []int{0}},
			},
		},
		{
			Name: "CrossedDeps",
			Blocks: []TestBlock{
				{name: "a", x: 0, y: 0},
				{name: "b", x: 2, y: 1},
				{name: "c", x: 4, y: 2},
				{name: "d", x: 6, y: 3},
				{name: "e", x: 8, y: 4},
			},
			Connections: []TestConnection{
				{from: 0, to: []int{1, 2, 3}},
				{from: 1, to: []int{3, 4}},
			},
		},
		{
			Name: "RoundTrip",
			Blocks: []TestBlock{
				{name: "a", x: 0, y: 0},
				{name: "b", x: 2, y: 1},
				{name: "c", x: 4, y: 2},
				{name: "d", x: 6, y: 3},
				{name: "e", x: 8, y: 4},
			},
			Connections: []TestConnection{
				{from: 0, to: []int{1, 2, 3}},
				{from: 1, to: []int{4}},
				{from: 2, to: []int{4}},
				{from: 3, to: []int{4}},
			},
		},
		{
			Name: "Two in same X",
			Blocks: []TestBlock{
				{name: "a", x: 2, y: 0},
				{name: "b", x: 2, y: 2},
				{name: "c", x: 4, y: 3},
			},
			Connections: []TestConnection{
				{from: 0, to: []int{2}},
				{from: 1, to: []int{2}},
			},
		},
		{
			Name: "Two in same X and an arrow in between",
			Blocks: []TestBlock{
				{name: "a", x: 0, y: 0},
				{name: "b", x: 3, y: 1},
				{name: "c", x: 3, y: 3},
				{name: "d", x: 5, y: 4},
			},
			Connections: []TestConnection{
				{from: 0, to: []int{1, 2, 3}},
				{from: 1, to: []int{3}},
				{from: 2, to: []int{3}},
			},
		},
		{
			Name: "Two in same X but one leaves space",
			Blocks: []TestBlock{
				{name: "a", x: 0, y: 0},
				{name: "b", x: 3, y: 1},
				{name: " c", x: 3, y: 3},
				{name: "d", x: 6, y: 4},
			},
			Connections: []TestConnection{
				{from: 1, to: []int{3}},
				{from: 2, to: []int{3}},
				{from: 0, to: []int{1, 2, 3}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			board := MakeBoard()
			for _, block := range tt.Blocks {
				err := board.AddBlock(block.name, block.name, block.x, block.y)
				a.NoError(err)
			}
			for _, connections := range tt.Connections {
				from := tt.Blocks[connections.from]
				for _, toI := range connections.to {
					to := tt.Blocks[toI]
					err := board.AddConnector(from.name, to.name)
					a.NoError(err)
				}
			}

			result, err := board.Render()
			a.NoError(err)
			expectTest(t, result)
		})
	}
}
