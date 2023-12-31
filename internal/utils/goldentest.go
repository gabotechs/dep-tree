package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func GoldenTest(t *testing.T, file string, content string) {
	dir := filepath.Dir(file)
	if !DirExists(dir) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}
	a := require.New(t)
	if FileExists(file) {
		expected, err := os.ReadFile(file)
		a.NoError(err)
		a.Equal(string(expected), content)
	} else {
		err := os.WriteFile(file, []byte(content), os.ModePerm)
		a.NoError(err)
	}
}
