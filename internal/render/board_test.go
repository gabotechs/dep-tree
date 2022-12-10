package render

import (
	"os"

	"github.com/stretchr/testify/require"

	"path"
	"testing"
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
	board := MakeBoard(BoardOptions{
		BlockSize: len("index.ts"),
	})
	blockA := "index.ts"
	blockB := "foo.ts"
	blockC := "bar.ts"

	_ = board.AddBlock(blockA, blockA, 0, 0)
	_ = board.AddBlock(blockB, blockB, 3, 4)
	_ = board.AddBlock(blockC, blockC, 5, 5)
	_ = board.AddDep(blockA, blockB)
	_ = board.AddDep(blockA, blockC)
	_ = board.AddDep(blockB, blockC)

	expectTest(t, board.Render())
}

func TestBoard_ReverseDeps(t *testing.T) {
	board := MakeBoard(BoardOptions{})
	one := "a"
	other := "b"

	_ = board.AddBlock(one, one, 0, 0)
	_ = board.AddBlock(other, other, 1, 1)
	_ = board.AddDep(other, one)

	expectTest(t, board.Render())
}

func TestBoard_CrossedDeps(t *testing.T) {
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
	_ = board.AddDep(blockA, blockB)
	_ = board.AddDep(blockA, blockC)
	_ = board.AddDep(blockA, blockD)
	_ = board.AddDep(blockB, blockD)
	_ = board.AddDep(blockB, blockE)

	expectTest(t, board.Render())
}
