package render

import (
	"fmt"

	"dep-tree/internal/render/graphics"
)

type Point struct {
	x int
	y int
}

type Block struct {
	Id       string
	Label    string
	Position Point
}

func (b *Block) Render(cells [][]graphics.CellStack) error {
	x := b.Position.x
	y := b.Position.y
	for i := 0; i < len(b.Label); i++ {
		idIndex := i
		if idIndex >= len(b.Label) || idIndex < 0 {
			// nothing here.
		} else {
			cells[y][x+i].PlaceChar(rune(b.Label[idIndex]))
		}
	}
	return nil
}

func (b *Board) AddBlock(id string, label string, column int, row int) error {
	if _, ok := b.blocks.Get(id); ok {
		return fmt.Errorf("block %s is already present", id)
	}

	x := column * b.options.Indent
	y := row

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
		Position: Point{x, y},
	})
	return nil
}
