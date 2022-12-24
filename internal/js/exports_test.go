package js

import (
	"context"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const exportsTestFolder = ".exports_test"

func TestParser_parseExports_IsCached(t *testing.T) {
	a := require.New(t)
	ctx := context.Background()
	file := path.Join(exportsTestFolder, "src", "index.js")
	parser, err := MakeJsParser(file)
	a.NoError(err)

	start := time.Now()
	ctx, _, err = parser.parseExports(ctx, file)
	a.NoError(err)
	nonCached := time.Since(start)

	start = time.Now()
	_, _, err = parser.parseExports(ctx, file)
	a.NoError(err)
	cached := time.Since(start)

	ratio := nonCached.Nanoseconds() / cached.Nanoseconds()
	a.Greater(ratio, int64(10))
}

func TestParser_parseExports(t *testing.T) {
	cwd, _ := os.Getwd()

	tests := []struct {
		Name     string
		File     string
		Expected map[string]string
	}{
		{
			Name: "test",
			File: path.Join(exportsTestFolder, "src", "index.js"),
			Expected: map[string]string{
				"Sorter":   path.Join(cwd, exportsTestFolder, "src", "utils", "sort.js"),
				"UnSorter": path.Join(cwd, exportsTestFolder, "src", "utils", "unsort.js"),
				"equals":   path.Join(cwd, exportsTestFolder, "src", "utils", "math", "equals.js"),
				"abs":      path.Join(cwd, exportsTestFolder, "src", "utils", "math", "index.js"),
				"sum":      path.Join(cwd, exportsTestFolder, "src", "utils", "math", "sum.js"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			parser, err := MakeJsParser(tt.File)
			a.NoError(err)
			_, exports, err := parser.parseExports(context.Background(), tt.File)
			a.NoError(err)
			a.Equal(tt.Expected, exports)
		})
	}
}
