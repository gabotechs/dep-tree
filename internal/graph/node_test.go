package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNodesIterator(t *testing.T) {
	a := require.New(t)
	iterator := NewNodesIterator([]*Node[int]{
		MakeNode("1", 1),
		MakeNode("2", 2),
		MakeNode("3", 3),
		MakeNode("4", 4),
	})

	a.Nil(iterator.Node())
	a.Equal(iterator.Next(), true)
	a.Equal(iterator.Node().ID(), hashCached("1"))
	a.Equal(iterator.Next(), true)
	a.Equal(iterator.Node().ID(), hashCached("2"))
	a.Equal(iterator.Next(), true)
	a.Equal(iterator.Node().ID(), hashCached("3"))
	a.Equal(iterator.Next(), true)
	a.Equal(iterator.Node().ID(), hashCached("4"))
	a.Equal(iterator.Next(), false)
	a.Nil(iterator.Node())
	iterator.Reset()
	a.Nil(iterator.Node())
	a.Equal(iterator.Next(), true)
	a.Equal(iterator.Node().ID(), hashCached("1"))
	a.Equal(iterator.Next(), true)
	a.Equal(iterator.Node().ID(), hashCached("2"))
	a.Equal(iterator.Next(), true)
	iterator.Reset()
	a.Nil(iterator.Node())
	a.Equal(iterator.Next(), true)
	a.Equal(iterator.Node().ID(), hashCached("1"))
	a.Equal(iterator.Next(), true)
}

func TestHashCached(t *testing.T) {
	a := require.New(t)
	a.Equal(hashCached("foo"), hash("foo"))
	a.Equal(hashCached("foo"), hashCached("foo"))
	a.Equal(MakeNode("foo", 1).ID(), hashCached("foo"))
}
