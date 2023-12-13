package config

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
)

const tmpFolder = "/tmp/dep-tree-check-tests"

func TestCheck(t *testing.T) {
	tests := []struct {
		Name     string
		Spec     [][]int
		Config   Config
		Failures []string
	}{
		{
			Name: "Simple",
			Spec: [][]int{
				{1, 2, 3},
				{2, 4},
				{3, 4},
				{4},
				{3},
			},
			Config: Config{
				Entrypoints: []string{"0"},
				WhiteList: map[string][]string{
					"4": {},
				},
				BlackList: map[string][]string{
					"0": {"3"},
				},
			},
			Failures: []string{
				"4 -> 3",
				"0 -> 3",
				"detected circular dependency: 3 -> 4 -> 3",
			},
		},
	}

	_ = os.MkdirAll(tmpFolder, os.ModePerm)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			ctx := context.Background()

			err := Check(ctx, func(ctx context.Context, s string) (context.Context, dep_tree.NodeParser[[]int], error) {
				return ctx, &dep_tree.TestParser{
					Start: s,
					Spec:  tt.Spec,
				}, nil
			}, &tt.Config) //nolint:gosec
			if tt.Failures != nil {
				msg := err.Error()
				failures := strings.Split(msg, "\n")
				failures = failures[1 : len(failures)-1]
				a.Equal(tt.Failures, failures)
			}
		})
	}
}
