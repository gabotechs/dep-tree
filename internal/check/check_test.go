package check

import (
	"strings"
	"testing"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/stretchr/testify/require"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		Name    string
		Spec    [][]int
		Config  *Config
		Failure string
	}{
		{
			Name: "Simple",
			Spec: [][]int{
				0: {1, 2, 3},
				1: {2, 4},
				2: {3, 4},
				3: {4},
				4: {3},
			},
			Config: &Config{
				Entrypoints: []string{"0"},
				WhiteList: map[string]WhiteListEntries{
					"4": {},
				},
				BlackList: map[string][]BlackListEntry{
					"0": {{To: "3"}},
				},
			},
			Failure: `
Check failed, the following dependencies are not allowed:
- 0 -> 3
- 4 -> 3

detected circular dependencies:
- 4 -> 3 -> 4`,
		},
		{
			Name: "With description",
			Spec: [][]int{
				0: {1, 2, 3},
				1: {2, 4},
				2: {3, 4},
				3: {4},
				4: {3},
			},
			Config: &Config{
				Entrypoints: []string{"0"},
				WhiteList: map[string]WhiteListEntries{
					"4": {Reason: "4 Should not be importing anything"},
				},
				BlackList: map[string][]BlackListEntry{
					"0": {{To: "3", Reason: "0 should not import 3"}},
				},
			},
			Failure: `
Check failed, the following dependencies are not allowed:
- 0 -> 3
  0 should not import 3
- 4 -> 3
  4 Should not be importing anything

detected circular dependencies:
- 4 -> 3 -> 4`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			err := Check[[]int](
				&graph.TestParser{Spec: tt.Spec},
				func(node *graph.Node[[]int]) string { return node.Id },
				tt.Config,
				nil,
			)
			if tt.Failure != "" {
				a.Equal(
					strings.TrimSpace(tt.Failure),
					strings.TrimSpace(err.Error()),
				)
			}
		})
	}
}
