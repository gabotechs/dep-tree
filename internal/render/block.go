package render

import (
	"fmt"

	"dep-tree/internal/render/graphics"
	"dep-tree/internal/vector"
)

const (
	cellType = "cellType"
	block    = "block"
	arrow    = "arrow"
)

type Block struct {
	Id       string
	Label    string
	Position vector.Vector
}

func (b *Block) Render(matrix *graphics.Matrix) error {
	x := b.Position.X
	y := b.Position.Y
	for i := 0; i < len(b.Label); i++ {
		idIndex := i
		if idIndex >= len(b.Label) || idIndex < 0 {
			// nothing here.
		} else {
			cell := matrix.Cell(vector.Vec(x+i, y))
			cell.PlaceChar(rune(b.Label[idIndex]))
			cell.Tag(cellType, block)
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
		return fmt.Errorf("block %s is already present", id)
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
		Position: vector.Vec(x, y),
	})
	return nil
}
