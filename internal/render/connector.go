package render

import (
	"fmt"

	"dep-tree/internal/render/graphics"
	"dep-tree/internal/utils"
	"dep-tree/internal/vector"
)

const (
	noCrossOwnership = "noCrossOwnership"
)

type Connector struct {
	from *Block
	to   *Block
}

func (c *Connector) Render(matrix *graphics.Matrix) error {
	reverseX := c.to.Position.X < c.from.Position.X
	reverseY := c.to.Position.Y < c.from.Position.Y

	dir := vector.Vec(utils.Bool2Int(!reverseX), utils.Bool2Int(!reverseY))
	// 1. If the line is going upwards, start at the end of the block.
	from := c.from.Position
	if reverseY {
		from.X += len(c.from.Label) - 1
	}

	// 2. start with just one vertical step.
	tracer := graphics.NewLineTracer(from)

	cur := tracer.MoveVertical(reverseY)
	cell := matrix.Cell(cur)
	if cell.Is(cellType, block) && cell.Is(cellType, arrow) {
		return fmt.Errorf("could not draw first vertical step on (%d, %d) because there is no space", cur.X, cur.Y)
	}

	// 3. move horizontally until no vertical collision is expected.
	for {
		collides, err := matrix.RayCastVertical(
			cur,
			map[string]func(string) bool{
				// if an arrow or a block is already present, then that is a hit.
				cellType: func(value string) bool {
					return utils.InArray(value, []string{block, arrow})
				},
				// if there is a line, belonging to another connector which claimed ownership, then hit.
				noCrossOwnership: func(value string) bool {
					return value != c.from.Id
				},
			},
			c.to.Position.Y-c.from.Position.Y,
		)
		if err != nil {
			return err
		} else if !collides {
			break
		}
		cur = tracer.MoveHorizontal(!reverseX)
		matrix.Cell(cur).Tag(noCrossOwnership, c.from.Id)
	}

	// 3. displacing vertically until aligned...
	for dir.Y*cur.Y < dir.Y*c.to.Position.Y {
		cur = tracer.MoveVertical(reverseY)
		matrix.Cell(cur).Tag(noCrossOwnership, c.from.Id)
	}

	// 4. moving horizontally until meeting target node...
	stopBefore := 1
	if dir.X < 0 {
		stopBefore = len(c.to.Label)
	}
	for dir.X*cur.X < dir.X*c.to.Position.X-stopBefore {
		cur = tracer.MoveHorizontal(reverseX)
	}
	err := tracer.Dump(matrix)
	if err != nil {
		return err
	}

	// 5. placing arrow in target node...
	cell = matrix.Cell(cur)
	cell.PlaceArrow(reverseX)
	cell.Tag(cellType, arrow)
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
