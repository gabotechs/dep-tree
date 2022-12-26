package board

import (
	"fmt"

	"dep-tree/internal/board/graphics"
	"dep-tree/internal/utils"
)

const (
	cellType = "cellType"

	blockChar  = "blockChar"
	blockSpace = "blockSpace"
	arrow      = "arrow"
)

type Block struct {
	Id       string
	Label    string
	Tags     map[string]string
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
				cell.PlaceChar(char).
					WithTag(cellType, blockChar).
					WithTags(b.Tags)
			} else {
				cell.PlaceEmpty().
					WithTag(cellType, blockSpace)
			}
		}
	}
	return nil
}

func (b *Board) AddBlock(block *Block) error {
	if _, ok := b.blocks.Get(block.Id); ok {
		return fmt.Errorf("blockChar %s is already present", block.Id)
	}

	newW := block.Position.X + len(block.Label)
	newH := block.Position.Y + 1

	if newW > b.w {
		b.w = newW
	}
	if newH > b.h {
		b.h = newH
	}

	b.blocks.Set(block.Id, block)
	return nil
}
