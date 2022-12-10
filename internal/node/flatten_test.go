package node

import (
	"github.com/stretchr/testify/require"

	"testing"
)

func TestNode_Flatten(t *testing.T) {
	a := require.New(t)

	node0 := MakeNode("0", 0)
	node1 := MakeNode("1", 0)
	node2 := MakeNode("2", 0)
	node3 := MakeNode("3", 0)

	node0.AddChild(node1)
	node0.AddChild(node2)
	node1.AddChild(node3)
	node2.AddChild(node3)

	result := node0.Flatten()
	a.Equal(result.Len(), 4)
	n0, _ := result.Get("0")
	a.Equal("0", n0.Id)
	n1, _ := result.Get("1")
	a.Equal("1", n1.Id)
	n2, _ := result.Get("2")
	a.Equal("2", n2.Id)
	n3, _ := result.Get("3")
	a.Equal("3", n3.Id)
}
