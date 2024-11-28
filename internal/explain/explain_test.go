package explain

import (
	"testing"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/stretchr/testify/require"
)

func TestExplain(t *testing.T) {
	tests := []struct {
		Name     string
		Spec     [][]int
		From     []string
		To       []string
		Expected []string
	}{
		{
			Name: "Simple",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {},
			},
			From:     []string{"0"},
			To:       []string{"1"},
			Expected: []string{"0 -> 1"},
		},
		{
			Name: "One to many",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {},
			},
			From:     []string{"0"},
			To:       []string{"1", "2"},
			Expected: []string{"0 -> 1"},
		},
		{
			Name: "Many to one",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {},
			},
			From:     []string{"0", "2"},
			To:       []string{"1"},
			Expected: []string{"0 -> 1"},
		},
		{
			Name: "Many to many",
			Spec: [][]int{
				0: {1},
				1: {2},
				2: {},
			},
			From:     []string{"0", "2"},
			To:       []string{"1", "2"},
			Expected: []string{"0 -> 1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			result, err := Explain[[]int](
				&graph.TestParser{Spec: tt.Spec},
				tt.From,
				tt.To,
				nil,
			)
			a.NoError(err)
			rendered := make([]string, len(result))
			for i, r := range result {
				rendered[i] = r[0].Id + " -> " + r[1].Id
			}
			a.Equal(tt.Expected, rendered)
		})
	}
}
