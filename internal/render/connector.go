package render

import (
	"fmt"

	"dep-tree/internal/render/graphics"
)

type Connector struct {
	from *Block
	to   *Block
}

func (c *Connector) Render(cells [][]graphics.CellStack) error {
	dir := Point{1, 1}
	if c.to.Position.x < c.from.Position.x {
		dir.x = -1
	}
	if c.to.Position.y < c.from.Position.y {
		dir.y = -1
	}

	cur := c.from.Position
	cur.y += dir.y
	if dir.y < 0 {
		cur.x += len(c.from.Label) - 1
	}
	for dir.y*cur.y < dir.y*c.to.Position.y {
		cells[cur.y][cur.x].DrawVerticalLine()
		cur.y += dir.y
	}
	cells[cur.y][cur.x].DrawJoint(dir.x > 0, dir.y < 0)

	stopBefore := 1
	if dir.x < 0 {
		stopBefore = len(c.to.Label)
	}

	for dir.x*cur.x < dir.x*c.to.Position.x-stopBefore {
		cur.x += dir.x
		cells[cur.y][cur.x].DrawHorizontalLine()
	}
	cells[cur.y][cur.x].PlaceArrow(dir.x < 0)
	return nil
}

func (b *Board) AddConnector(from string, to string) error {
	var fromBlock *Block
	var toBlock *Block
	if block, ok := b.blocks.Get(from); ok {
		fromBlock = block
	} else {
		return fmt.Errorf("block with Id %s not found", from)
	}
	if block, ok := b.blocks.Get(to); ok {
		toBlock = block
	} else {
		return fmt.Errorf("block with Id %s not found", to)
	}

	key := from + " -> " + to
	if _, ok := b.connectors.Get(key); ok {
		return fmt.Errorf("connector from %s to %s already present", from, to)
	}
	b.connectors.Set(key, &Connector{
		from: fromBlock,
		to:   toBlock,
	})
	return nil
}
