package render

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
	fullPath := path.Join(testPath, t.Name()+".txt")
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

func TestBoard_SimpleDeps(t *testing.T) {
	a := require.New(t)
	board := MakeBoard(BoardOptions{})
	blockA := "index.ts"
	blockB := "foo.ts"
	blockC := "bar.ts"

	_ = board.AddBlock(blockA, blockA, 0, 0)
	_ = board.AddBlock(blockB, blockB, 3, 4)
	_ = board.AddBlock(blockC, blockC, 5, 5)
	_ = board.AddConnector(blockA, blockB)
	_ = board.AddConnector(blockA, blockC)
	_ = board.AddConnector(blockB, blockC)

	result, err := board.Render()
	a.NoError(err)
	expectTest(t, result)
}

func TestBoard_ReverseDeps(t *testing.T) {
	a := require.New(t)
	board := MakeBoard(BoardOptions{})
	one := "a"
	other := "b"

	_ = board.AddBlock(one, one, 0, 0)
	_ = board.AddBlock(other, other, 1, 1)
	_ = board.AddConnector(other, one)

	result, err := board.Render()
	a.NoError(err)
	expectTest(t, result)
}

func TestBoard_CrossedDeps(t *testing.T) {
	a := require.New(t)
	board := MakeBoard(BoardOptions{})
	blockA := "a"
	blockB := "b"
	blockC := "c"
	blockD := "d"
	blockE := "e"

	_ = board.AddBlock(blockA, blockA, 0, 0)
	_ = board.AddBlock(blockB, blockB, 1, 1)
	_ = board.AddBlock(blockC, blockC, 2, 2)
	_ = board.AddBlock(blockD, blockD, 3, 3)
	_ = board.AddBlock(blockE, blockE, 4, 4)
	_ = board.AddConnector(blockA, blockB)
	_ = board.AddConnector(blockA, blockC)
	_ = board.AddConnector(blockA, blockD)
	_ = board.AddConnector(blockB, blockD)
	_ = board.AddConnector(blockB, blockE)

	result, err := board.Render()
	a.NoError(err)
	expectTest(t, result)
}

func TestBoard_RoundTripSolid(t *testing.T) {
	a := require.New(t)
	board := MakeBoard(BoardOptions{})
	blockA := "a"
	blockB := "b"
	blockC := "c"
	blockD := "d"
	blockE := "e"

	_ = board.AddBlock(blockA, blockA, 0, 0)
	_ = board.AddBlock(blockB, blockB, 2, 1)
	_ = board.AddBlock(blockC, blockC, 2, 3)
	_ = board.AddBlock(blockD, blockD, 2, 5)
	_ = board.AddBlock(blockE, blockE, 4, 6)

	_ = board.AddConnector(blockA, blockB)
	_ = board.AddConnector(blockA, blockC)
	_ = board.AddConnector(blockA, blockD)

	_ = board.AddConnector(blockC, blockE)
	_ = board.AddConnector(blockD, blockE)
	_ = board.AddConnector(blockB, blockE)

	result, err := board.Render()
	a.NoError(err)
	expectTest(t, result)
}
