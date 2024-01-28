package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/graph/topo"
)

func TestGraph_AddNode(t *testing.T) {
	a := require.New(t)
	g := NewGraph[int]()
	g.AddNode(MakeNode[int]("0", 0))
	g.AddNode(MakeNode[int]("1", 1))
	g.AddNode(MakeNode[int]("1", 1))
	a.Equal(2, len(g.AllNodes()))
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

func TestGraph_AddFromToEdge(t *testing.T) {
	a := require.New(t)
	g := NewGraph[int]()
	g.AddNode(MakeNode[int]("0", 0))
	g.AddNode(MakeNode[int]("1", 1))
	err := g.AddFromToEdge("0", "1")
	a.NoError(err)
	err = g.AddFromToEdge("0", "2")
	a.Error(err)
	err = g.AddFromToEdge("2", "0")
	a.Error(err)
	err = g.AddFromToEdge("2", "3")
	a.Error(err)

	a.Equal(1, len(g.FromId("0")))
	a.Equal(1, g.FromId("0")[0].Data)
	a.Equal(0, len(g.ToId("0")))

	a.Equal(1, len(g.ToId("1")))
	a.Equal(0, g.ToId("1")[0].Data)
	a.Equal(0, len(g.FromId("1")))
}

func TestGraph_Cycles(t *testing.T) {
	a := require.New(t)
	g := NewGraph[int]()

	node0 := MakeNode[int]("0", 0)
	node1 := MakeNode[int]("1", 1)
	node2 := MakeNode[int]("2", 2)
	node3 := MakeNode[int]("3", 3)

	g.AddNode(node0)
	g.AddNode(node1)
	g.AddNode(node2)
	g.AddNode(node3)
	err := g.AddFromToEdge("0", "1")
	a.NoError(err)
	err = g.AddFromToEdge("1", "2")
	a.NoError(err)
	err = g.AddFromToEdge("1", "0")
	a.NoError(err)
	err = g.AddFromToEdge("2", "3")
	a.NoError(err)
	err = g.AddFromToEdge("3", "0")
	a.NoError(err)
	cycles := topo.DirectedCyclesIn(g)
	a.Equal(len(cycles), 2)
}

func TestGraph_GetNodesWithoutParents(t *testing.T) {
	a := require.New(t)
	g := NewGraph[int]()

	node0 := MakeNode[int]("0", 0)
	node1 := MakeNode[int]("1", 1)

	g.AddNode(node0)
	g.AddNode(node1)
	err := g.AddFromToEdge("0", "1")
	a.NoError(err)

	nodes := g.GetNodesWithoutParents()
	a.Equal(1, len(nodes))
	a.Equal("0", nodes[0].Id)

	err = g.AddFromToEdge("1", "0")
	a.NoError(err)
	nodes = g.GetNodesWithoutParents()
	a.Equal(0, len(nodes))
}
