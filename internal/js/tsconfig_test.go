package js

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTsConfig_ResolveFromPaths(t *testing.T) {
	tests := []struct {
		Name     string
		TsConfig TsConfig
		Import   string
		Result   []string
	}{
		{
			Name: "without globstar",
			TsConfig: TsConfig{
				path: "/foo",
				CompilerOptions: CompilerOptions{
					BaseUrl: "./bar/baz",
					Paths:   map[string][]string{"@/": {"./src/"}},
				},
			},
			Import: "@/a/b/c",
			Result: []string{"/foo/bar/baz/src/a/b/c"},
		},
		{
			Name: "without globstar (2)",
			TsConfig: TsConfig{
				path: "/foo",
				CompilerOptions: CompilerOptions{
					BaseUrl: "./bar/baz",
					Paths:   map[string][]string{"@Environment": {"./src/environments/environments.ts"}},
				},
			},
			Import: "@Environment",
			Result: []string{"/foo/bar/baz/src/environments/environments.ts"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			result := tt.TsConfig.ResolveFromPaths(tt.Import)
			a.Equal(tt.Result, result)
		})
	}
}
