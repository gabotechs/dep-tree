package golang

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoMod(t *testing.T) {
	tests := []struct {
		Name     string
		Expected GoMod
	}{
		{
			Name: "../../go.mod",
			Expected: GoMod{
				Module: "github.com/gabotechs/dep-tree",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			result, err := ParseGoMod(tt.Name)
			a.NoError(err)
			a.Equal(tt.Expected, *result)
		})
	}
}
