package board

import (
	"fmt"

	"github.com/gabotechs/dep-tree/internal/board/graphics"
	"github.com/gabotechs/dep-tree/internal/utils"
)

type Connector struct {
	from *Block
	to   *Block
	tags map[string]string
}

func checkCollision(
	matrix *graphics.Matrix,
	cursor utils.Vector,
	from *Block,
	to *Block,
) (bool, error) {
	return matrix.RayCastVertical(
		cursor,
		map[string]func(string) bool{
			// if an arrow or a blockChar is already present, then that is a hit.
			cellType: func(value string) bool {
				return utils.InArray(value, []string{blockChar, arrow})
			},
		},
		to.Position.Y-from.Position.Y,
	)
}

// Render This function is far less clear if factored out.
//
//nolint:gocyclo
func (c *Connector) Render(matrix *graphics.Matrix) error {
	reverseX := c.to.Position.X <= c.from.Position.X
	reverseY := c.to.Position.Y < c.from.Position.Y

	// 1. If the line is going upwards, start at the end of the blockChar.
	from := c.from.Position
	if reverseY {
		from.X += len(c.from.Label) - 1
	} else {
		from.X += utils.PrefixN(c.from.Label, ' ')
	}
	// 2. start with just one vertical step.
	tracer := graphics.
		NewLineTracer(from).
		WithTags(c.tags)

	var cur utils.Vector
	if reverseY {
		cur = tracer.MoveHorizontal(false)
		if cur.X >= matrix.W() {
			matrix.ExpandRight(1)
		}
	} else {
		cur = tracer.MoveVertical(false)
	}
	cell := matrix.Cell(cur)
	if cell == nil || cell.Is(cellType, blockChar) { // || cell.Is(cellType, arrow)
		return fmt.Errorf("could not draw first vertical step on (%d, %d) because there is no space", cur.X, cur.Y)
	}

	// 3. move horizontally until no vertical collision is expected.
	for {
		collides, err := checkCollision(matrix, cur, c.from, c.to)
		if err != nil {
			return err
		} else if !collides {
			break
		}

		cur = tracer.MoveHorizontal(!reverseX)
		cell = matrix.Cell(cur)
		if cell == nil && cur.X >= matrix.W() {
			// While moving horizontally, the end was reached, expanding...
			matrix.ExpandRight(1)
			cell = matrix.Cell(cur)
		}
		if cell == nil {
			return fmt.Errorf("moved to invalid position (%d, %d) while tracing horizontal line. This error should not be reachable", cur.X, cur.Y)
		}
	}

	// 4. displacing vertically until aligned...
	for cur.Y != c.to.Position.Y && cur.Y >= 0 && cur.Y < matrix.H() {
		cur = tracer.MoveVertical(cur.Y > c.to.Position.Y)
	}

	// 5. moving horizontally until meeting target node...
	for cur.X != c.to.Position.X && cur.X >= 0 && cur.X < matrix.W() {
		next := matrix.Cell(utils.Vec(cur.X+utils.Bool2Int(!reverseX), cur.Y))
		if next != nil && (next.Is(cellType, blockChar) || next.Is(cellType, blockSpace)) {
			break
		}
		cur = tracer.MoveHorizontal(cur.X > c.to.Position.X)
	}
	err := tracer.Dump(matrix)
	if err != nil {
		return err
	}

	// 6. placing arrow in target node...
	cell = matrix.Cell(cur)
	if cell != nil {
		cell.PlaceArrow(reverseX).
			WithTag(cellType, arrow).
			WithTags(c.tags)
	}
	return nil
}

func (b *Board) AddConnector(
	from string,
	to string,
	tags map[string]string,
) error {
	var fromBlock *Block
	var toBlock *Block
	if block, ok := b.blocks.Get(from); ok {
		fromBlock = block
	} else {
		return fmt.Errorf("blockChar with Id %s not found", from)
	}
	if block, ok := b.blocks.Get(to); ok {
		toBlock = block
	} else {
		return fmt.Errorf("blockChar with Id %s not found", to)
	}

	key := from + " -> " + to
	if _, ok := b.connectors.Get(key); ok {
		return fmt.Errorf("connector from %s to %s already present", from, to)
	}
	b.connectors.Set(key, &Connector{
		from: fromBlock,
		to:   toBlock,
		tags: tags,
	})
	return nil
}
