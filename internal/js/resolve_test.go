package js

import (
	"context"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const resolverTestFolder = ".resolve_test"

func TestParser_ResolvePath(t *testing.T) {
	absPath, _ := filepath.Abs(resolverTestFolder)

	tests := []struct {
		Name       string
		Unresolved string
		Cwd        string
		Resolved   string
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
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			parser, err := MakeJsParser(tt.Cwd)
			a.NoError(err)
			_, resolved, err := parser.ResolvePath(context.Background(), tt.Unresolved, tt.Cwd)
			a.NoError(err)
			a.Equal(tt.Resolved, resolved)
		})
	}
}
