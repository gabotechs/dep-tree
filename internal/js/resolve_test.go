package js

import (
	"context"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const resolverTestFolder = ".resolve_test"

func TestParser_ResolvePath_IsCached(t *testing.T) {
	a := require.New(t)
	ctx := context.Background()
	parser, err := MakeJsParser(resolverTestFolder)
	a.NoError(err)

	start := time.Now()
	ctx, _, err = parser.ResolvePath(ctx, path.Join(resolverTestFolder, "src", "foo.ts"), resolverTestFolder)
	a.NoError(err)
	nonCached := time.Since(start)

	start = time.Now()
	_, _, err = parser.ResolvePath(ctx, path.Join(resolverTestFolder, "src", "foo.ts"), resolverTestFolder)
	a.NoError(err)
	cached := time.Since(start)

	ratio := nonCached.Nanoseconds() / cached.Nanoseconds()
	a.Greater(ratio, int64(5))
}

func TestParser_ResolvePath(t *testing.T) {
	absPath, _ := filepath.Abs(resolverTestFolder)

	tests := []struct {
		Name          string
		Unresolved    string
		Cwd           string
		Resolved      string
		ExpectedError string
	}{
		{
			Name:       "from relative",
			Cwd:        path.Join(resolverTestFolder, "src", "utils"),
			Unresolved: "../foo",
			Resolved:   path.Join(absPath, "src", "foo.ts"),
		},
		{
			Name:       "from baseUrl",
			Cwd:        path.Join(resolverTestFolder, "src"),
			Unresolved: "foo",
			Resolved:   path.Join(absPath, "src", "foo.ts"),
		},
		{
			Name:       "from paths override",
			Cwd:        path.Join(resolverTestFolder, "src"),
			Unresolved: "@utils/sum",
			Resolved:   path.Join(absPath, "src", "utils", "sum.ts"),
		},
		{
			Name:       "from paths override with glob pattern",
			Cwd:        path.Join(resolverTestFolder, "src"),
			Unresolved: "@/helpers/diff",
			Resolved:   path.Join(absPath, "src", "helpers", "diff.ts"),
		},
		{
			Name:          "Does not resolve invalid import",
			Cwd:           path.Join(resolverTestFolder, "src"),
			Unresolved:    "bar",
			ExpectedError: "import 'bar' cannot be resolved",
		},
		{
			Name:          "Does not resolve invalid relative import",
			Cwd:           path.Join(resolverTestFolder, "src", "utils"),
			Unresolved:    "./foo",
			ExpectedError: "could not perform relative import for './foo' because the file or dir was not found",
		},
		{
			Name:          "Does not resolve invalid relative import",
			Cwd:           resolverTestFolder,
			Unresolved:    path.Join("src", "utils", "foo"),
			ExpectedError: "import 'src/utils/foo' cannot be resolved",
		},
		{
			Name:          "Does not resolve invalid path override import",
			Cwd:           path.Join(resolverTestFolder, "src"),
			Unresolved:    "@/helpers/bar",
			ExpectedError: "import '@/helpers/bar' was matched to path '@/helpers/' in tscofing's paths option, but the resolved path did not match an existing file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			parser, err := MakeJsParser(tt.Cwd)
			a.NoError(err)
			_, resolved, err := parser.ResolvePath(context.Background(), tt.Unresolved, tt.Cwd)
			if tt.ExpectedError != "" {
				a.ErrorContains(err, tt.ExpectedError)
			} else {
				a.NoError(err)
				a.Equal(tt.Resolved, resolved)
			}
		})
	}
}
