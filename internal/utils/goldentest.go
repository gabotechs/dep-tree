package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func GoldenTest(t *testing.T, file string, content string) {
	dir := filepath.Dir(file)
	if !DirExists(dir) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}
	a := require.New(t)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	if FileExists(file) {
		expected, err := os.ReadFile(file)
		a.NoError(err)
		a.Equal(string(expected), content)
	} else {
		err := os.WriteFile(file, []byte(content), os.ModePerm)
		a.NoError(err)
	}
}
