package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveOverlap(t *testing.T) {
	tests := []struct {
		Name     string
		A        []string
		B        []string
		Expected []string
	}{
		{
			Name:     "No overlap",
			A:        []string{"a", "b", "c"},
			B:        []string{"d", "e", "f"},
			Expected: []string{"a", "b", "c"},
		},
		{
			Name:     "Partial overlap",
			A:        []string{"a", "b", "c"},
			B:        []string{"b", "c", "d"},
			Expected: []string{"a"},
		},
		{
			Name:     "Full overlap",
			A:        []string{"a", "b", "c"},
			B:        []string{"a", "b", "c"},
			Expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Expected, RemoveOverlap(tt.A, tt.B))

		})
	}
}
