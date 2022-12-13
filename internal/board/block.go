package board

import (
	"fmt"

	"dep-tree/internal/board/graphics"
	"dep-tree/internal/utils"
)

const (
	cellType   = "cellType"
	blockChar  = "blockChar"
	blockSpace = "blockSpace"
	arrow      = "arrow"
)

type Block struct {
	Id       string
	Label    string
	Position utils.Vector
}

func (b *Block) Render(matrix *graphics.Matrix) error {
	x := b.Position.X
	y := b.Position.Y
	for i := 0; i < len(b.Label); i++ {
		idIndex := i
		if idIndex >= len(b.Label) || idIndex < 0 {
			// nothing here.
		} else {
			cell := matrix.Cell(utils.Vec(x+i, y))
			if cell == nil {
				return fmt.Errorf("tried to render in invalid cell (%d, %d)", x+i, y)
			}
			char := rune(b.Label[idIndex])
			if char != ' ' {
				cell.PlaceChar(char)
				cell.Tag(cellType, blockChar)
			} else {
				cell.Tag(cellType, blockSpace)
			}
		}
	}
	return nil
}

func (b *Board) AddBlock(
	id string,
	label string,
	x int,
	y int,
) error {
	if _, ok := b.blocks.Get(id); ok {
		return fmt.Errorf("blockChar %s is already present", id)
	}

	newW := x + len(label)
	newH := y + 1

	if newW > b.w {
		b.w = newW
	}
	if newH > b.h {
		b.h = newH
	}

	b.blocks.Set(id, &Block{
		Id:       id,
		Label:    label,
		Position: utils.Vec(x, y),
	})
	return nil
}
