package graphics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCell_WithTag(t *testing.T) {
	a := require.New(t)
	c := TaggedCell{}

	a.Equal(false, c.Is("key", "foo"))
	c.WithTag("key", "bar")
	a.Equal(false, c.Is("key", "foo"))
	a.Equal(true, c.Is("key", "bar"))
	c.WithTags(map[string]string{
		"key":      "foo",
		"otherKey": "bar",
	})
	a.Equal(true, c.Is("key", "foo"))
	a.Equal(true, c.Is("otherKey", "bar"))
}
