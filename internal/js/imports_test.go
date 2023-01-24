package js

import (
	"context"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const importsTestFolder = ".imports_test"

func TestParser_parseImports_IsCached(t *testing.T) {
	a := require.New(t)
	ctx := context.Background()
	file := path.Join(importsTestFolder, "index.ts")
	p, err := MakeJsParser(file)
	a.NoError(err)
	parser := p.(*Parser)

	start := time.Now()
	ctx, _, err = parser.parseImports(ctx, file)
	a.NoError(err)
	nonCached := time.Since(start)

	start = time.Now()
	_, _, err = parser.parseImports(ctx, file)
	a.NoError(err)
	cached := time.Since(start)

	ratio := nonCached.Nanoseconds() / cached.Nanoseconds()
	a.Greater(ratio, int64(10))
}

func TestParser_parseImports(t *testing.T) {
	wd, _ := os.Getwd()

	tests := []struct {
		Name           string
		File           string
		Expected       map[string][]string
		ExpectedErrors []string
	}{
		{
			Name: "test 1",
			File: path.Join(importsTestFolder, "index.ts"),
			Expected: map[string][]string{
				path.Join(wd, importsTestFolder, "2", "2.ts"):      {"a", "b"},
				path.Join(wd, importsTestFolder, "2", "index.ts"):  {"*"},
				path.Join(wd, importsTestFolder, "1", "a", "a.ts"): {"*"},
			},
			ExpectedErrors: []string{
				"could not perform relative import for './unexisting'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			p, err := MakeJsParser(tt.File)
			a.NoError(err)
			parser := p.(*Parser)
			_, results, err := parser.parseImports(context.Background(), tt.File)
			a.NoError(err)
			for expectedPath, expectedNames := range tt.Expected {
				resultNames, ok := results.Imports.Get(expectedPath)
				a.Equal(true, ok)
				a.Equal(expectedNames, resultNames)
			}

			a.Equal(len(tt.ExpectedErrors), len(results.Errors))
			if results.Errors != nil {
				for i, err := range results.Errors {
					a.ErrorContains(err, tt.ExpectedErrors[i])
				}
			}
		})
	}
}
