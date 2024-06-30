package js

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTsConfig_ResolveFromPaths(t *testing.T) {
	tests := []struct {
		Name     string
		basePath string
		Paths    map[string][]string
		Import   string
		Result   []string
	}{
		{
			Name:     "without globstar",
			basePath: "/foo/bar",
			Paths:    map[string][]string{"@/": {"./src/"}},
			Import:   "@/a/b/c",
			Result:   []string{"/foo/bar/src/a/b/c"},
		},
		{
			Name:     "without globstar (2)",
			basePath: "/foo/bar",
			Paths:    map[string][]string{"@Environment": {"./src/environments/environments.ts"}},
			Import:   "@Environment",
			Result:   []string{"/foo/bar/src/environments/environments.ts"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			tsConfig := TsConfig{
				path: tt.basePath,
				CompilerOptions: CompilerOptions{
					Paths: tt.Paths,
				},
			}
			result := tsConfig.ResolveFromPaths(tt.Import)
			a.Equal(tt.Result, result)
		})
	}
}
