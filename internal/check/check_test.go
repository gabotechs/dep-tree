package check

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
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
				"detected circular dependency: 3 -> 4 -> 3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			err := Check[[]int](&dep_tree.TestParser{Spec: tt.Spec}, tt.Config)
			if tt.Failures != nil {
				msg := err.Error()
				failures := strings.Split(msg, "\n")
				failures = failures[1:]
				a.Equal(tt.Failures, failures)
			}
		})
	}
}
