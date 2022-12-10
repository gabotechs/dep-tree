package node

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNode_Level(t *testing.T) {
	a := require.New(t)

	node0 := MakeNode("0", 0)
	node1 := MakeNode("1", 0)
	node2 := MakeNode("2", 0)
	node3 := MakeNode("3", 0)

	node0.AddChild(node1, node2)
	node1.AddChild(node3)
	node2.AddChild(node3)

	ctx := context.Background()
	ctx, lvl0 := node0.Level(ctx, "0")
	ctx, lvl1 := node1.Level(ctx, "0")
	ctx, lvl2 := node2.Level(ctx, "0")
	_, lvl3 := node3.Level(ctx, "0")

	a.Equal(0, lvl0)
	a.Equal(1, lvl1)
	a.Equal(1, lvl2)
	a.Equal(2, lvl3)
}

func TestNode_Level_Circular(t *testing.T) {
	a := require.New(t)

	node0 := MakeNode("0", 0)
	node1 := MakeNode("1", 0)
	node2 := MakeNode("2", 0)
	node3 := MakeNode("3", 0)
	node4 := MakeNode("4", 0)

	node0.AddChild(node1, node2, node3)
	node1.AddChild(node2, node4)
	node2.AddChild(node3, node4)
	node3.AddChild(node4)
	node4.AddChild(node3)

	ctx := context.Background()
	ctx, lvl0 := node0.Level(ctx, "0")
	ctx, lvl1 := node1.Level(ctx, "0")
	ctx, lvl2 := node2.Level(ctx, "0")
	ctx, lvl3 := node3.Level(ctx, "0")
	_, lvl4 := node4.Level(ctx, "0")

	a.Equal(0, lvl0)
	a.Equal(1, lvl1)
	a.Equal(2, lvl2)
	a.Equal(4, lvl3)
	a.Equal(3, lvl4)
}
