package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	a := require.New(t)

	tests := []struct {
		Name     string
		Pattern  string
		Path     string
		Expected bool
	}{
		{
			Name:     "1",
			Pattern:  "**",
			Path:     "/this/matches/anything.txt",
			Expected: true,
		},
		{
			Name:     "2",
			Pattern:  "*",
			Path:     "/this/does/not/match.txt",
			Expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			pass, err := GlobstarMatch(tt.Pattern, tt.Path)
			a.NoError(err)
			a.Equal(tt.Expected, pass)
		})
	}
}
