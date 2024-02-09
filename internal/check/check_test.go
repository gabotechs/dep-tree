package check

import (
	"strings"
	"testing"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/stretchr/testify/require"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		Name     string
		Spec     [][]int
		Config   *Config
		Failures []string
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
				WhiteList: map[string][]string{
					"4": {},
				},
				BlackList: map[string][]string{
					"0": {"3"},
				},
			},
			Failures: []string{
				"0 -> 3",
				"4 -> 3",
				"detected circular dependency: 4 -> 3 -> 4",
			},
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
			if tt.Failures != nil {
				msg := err.Error()
				failures := strings.Split(msg, "\n")
				failures = failures[1:]
				a.Equal(tt.Failures, failures)
			}
		})
	}
}
