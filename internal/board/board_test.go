package board

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const RebuildTestsEnv = "REBUILD_TESTS"

const testPath = ".board_test"

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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
		Name          string
		Blocks        []TestBlock
		Connections   []TestConnection
		ExpectedError string
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
			Name: "Cannot draw line when one is just above",
			Blocks: []TestBlock{
				{name: "foo.ts", x: 3, y: 4},
				{name: "bar.ts", x: 3, y: 5},
			},
			Connections: []TestConnection{
				{from: 0, to: []int{1}},
			},
			ExpectedError: "could not draw first vertical step on (3, 5) because there is no space",
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
			Name: "it should increase the board's size if necessary",
			Blocks: []TestBlock{
				{name: "long", x: 0, y: 0},
				{name: "b", x: 2, y: 1},
			},
			Connections: []TestConnection{
				{from: 1, to: []int{0}},
			},
		},
		{
			Name: "it should increase the board's size if necessary 2",
			Blocks: []TestBlock{
				{name: "a", x: 1, y: 0},
				{name: "bb", x: 0, y: 2},
			},
			Connections: []TestConnection{
				{from: 1, to: []int{0}},
			},
		},
		{
			Name: "it should increase the board's size if necessary 3",
			Blocks: []TestBlock{
				{name: "a", x: 1, y: 0},
				{name: "bb", x: 0, y: 2},
			},
			Connections: []TestConnection{
				{from: 0, to: []int{1}},
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
		{
			Name: "Reverse dep with first node very long",
			Blocks: []TestBlock{
				{name: "some-really-long-file", x: 0, y: 0},
				{name: "b", x: 3, y: 1},
			},
			Connections: []TestConnection{
				{from: 0, to: []int{1}},
				{from: 1, to: []int{0}},
			},
		},
		{
			Name: "Reverse dep with exactly same length",
			Blocks: []TestBlock{
				{name: "aaa", x: 0, y: 0},
				{name: "b", x: 2, y: 1},
			},
			Connections: []TestConnection{
				{from: 0, to: []int{1}},
				{from: 1, to: []int{0}},
			},
		},
		{
			Name: "Cyclic deps",
			Blocks: []TestBlock{
				{name: "a", x: 0, y: 0},
				{name: "b", x: 2, y: 1},
				{name: "c", x: 4, y: 2},
			},
			Connections: []TestConnection{
				{from: 0, to: []int{1}},
				{from: 1, to: []int{2}},
				{from: 2, to: []int{1}},
			},
		},
		{
			Name: "Does not hang",
			Blocks: []TestBlock{
				{name: "a", x: 0, y: 0},
				{name: "b", x: 2, y: 1},
				{name: "c", x: 4, y: 2},
				{name: "d", x: 8, y: 3},
				{name: " e", x: 8, y: 4},
			},
			Connections: []TestConnection{
				{from: 0, to: []int{1, 2, 3}},
				{from: 1, to: []int{2, 4}},
				{from: 2, to: []int{3, 4}},
				{from: 3, to: []int{4}},
				{from: 4, to: []int{3}},
			},
		},
		{
			Name: "Reverse deps are drawn on the right",
			Blocks: []TestBlock{
				{name: " ../../../i18n/strings.ts", x: 8, y: 0},
				{name: "  ../../../services/ucloud/ackWebsocket/baseSocket.ts", x: 8, y: 1},
				{name: "   ../../../services/user.ts", x: 8, y: 2},
			},
			Connections: []TestConnection{
				{from: 2, to: []int{0}},
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
			if tt.ExpectedError == "" {
				a.NoError(err)
				_ = os.Mkdir(testPath, os.ModePerm)
				fullPath := path.Join(testPath, path.Base(t.Name())+".txt")
				print(result)
				if fileExists(fullPath) && os.Getenv(RebuildTestsEnv) != "true" {
					expected, err := os.ReadFile(fullPath)
					a.NoError(err)
					a.Equal(string(expected), result)
				} else {
					err := os.WriteFile(fullPath, []byte(result), os.ModePerm)
					a.NoError(err)
				}
			} else {
				a.ErrorContains(err, tt.ExpectedError)
			}
		})
	}
}
