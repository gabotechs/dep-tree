package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGraph_AddNode(t *testing.T) {
	a := require.New(t)
	g := NewGraph[int]()
	g.AddNode(MakeNode[int]("0", 0))
	g.AddNode(MakeNode[int]("1", 1))
	g.AddNode(MakeNode[int]("1", 1))
	a.Equal(2, len(g.Nodes()))
	a.Equal(true, g.Has("0"))
	a.Equal(true, g.Has("1"))
	a.Equal(false, g.Has("2"))

	a.Equal(g.Get("0"), MakeNode[int]("0", 0))
	a.Equal(g.Get("1"), MakeNode[int]("1", 1))
	a.Nil(g.Get("2"))

	a.Equal(true, g.Has("0"))
	a.Equal(true, g.Has("1"))
	a.Equal(false, g.Has("2"))
}

func TestGraph_AddChild(t *testing.T) {
	a := require.New(t)
	g := NewGraph[int]()
	g.AddNode(MakeNode[int]("0", 0))
	g.AddNode(MakeNode[int]("1", 1))
	err := g.AddChild("0", "1")
	a.NoError(err)
	err = g.AddChild("0", "2")
	a.Error(err)
	err = g.AddChild("2", "0")
	a.Error(err)
	err = g.AddChild("2", "3")
	a.Error(err)

	a.Equal(1, len(g.Children("0")))
	a.Equal(1, g.Children("0")[0].Data)
	a.Equal(0, len(g.Parents("0")))

	a.Equal(1, len(g.Parents("1")))
	a.Equal(0, g.Parents("1")[0].Data)
	a.Equal(0, len(g.Children("1")))
}
