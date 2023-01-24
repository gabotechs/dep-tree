package dep_tree

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"dep-tree/internal/config"
)

func TestValidateGraph(t *testing.T) {
	tests := []struct {
		Name     string
		Spec     [][]int
		Config   config.Config
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
			Config: config.Config{
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			testParser := TestGraph{
				Start: "0",
				Spec:  tt.Spec,
			}

			ctx := context.Background()

			_, dt, err := NewDepTree[[]int](ctx, &testParser)
			a.NoError(err)

			err = dt.Validate(&tt.Config)
			if err == nil {
				a.Equal(tt.Failures, 0)
			} else {
				a.Equal(tt.Failures, strings.Split(err.Error(), "\n")[1:])
			}
		})
	}
}
