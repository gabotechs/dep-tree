package render

import (
	"dep-tree/internal/utils"
	"fmt"
)

type Point struct {
	x int
	y int
}

type Board struct {
	w               int
	h               int
	indent          int
	blockSize       int
	blockAlignRight bool
	elements        [][]CellStack
	blocks          map[string]Point
}

type BoardOptions struct {
	Indent          int
	BlockSize       int
	BlockAlignRight bool
}

func MakeBoard(options BoardOptions) *Board {
	elements := make([][]CellStack, 0)
	return &Board{
		w:         0,
		h:         0,
		blockSize: utils.Clamp(6, options.BlockSize, 50),
		indent:    utils.Clamp(2, options.Indent, 8),
		elements:  elements,
		blocks:    map[string]Point{},
	}
}

func (b *Board) resize(x int, y int) {
	if !(x >= b.w || y >= b.h) {
		return
	}
	for i := 0; i <= y; i++ {
		if i >= len(b.elements) {
			b.elements = append(b.elements, make([]CellStack, 0))
		}
		for j := 0; j <= x; j++ {
			if j >= len(b.elements[i]) {
				b.elements[i] = append(b.elements[i], CellStack{cells: []Cell{}})
			}
		}
	}
	b.w = x + 1
	b.h = y + 1
}

func (b *Board) AddBlock(id string, display string, level int, index int) error {
	if _, ok := b.blocks[id]; ok {
		return fmt.Errorf("block %s is already present", id)
	}

	y := index
	x := b.indent * level
	b.resize(x+b.blockSize-1, y)
	for i := 0; i < b.blockSize; i++ {
		idIndex := i
		if b.blockAlignRight {
			idIndex -= b.blockSize - len(display)
		}
		if idIndex >= len(display) || idIndex < 0 {
			b.elements[y][x+i].PlaceChar(' ')
		} else {
			b.elements[y][x+i].PlaceChar(rune(display[idIndex]))
		}
	}
	b.blocks[id] = Point{x, y}
	return nil
}

func (b *Board) AddDep(from string, to string) error {
	var fromXY Point
	var toXY Point
	if coords, ok := b.blocks[from]; ok {
		fromXY = coords
	} else {
		return fmt.Errorf("block with Id %s not found", from)
	}
	if coords, ok := b.blocks[to]; ok {
		toXY = coords
	} else {
		return fmt.Errorf("block with Id %s not found", to)
	}

	dir := Point{1, 1}
	if toXY.x < fromXY.x {
		dir.x = -1
	}
	if toXY.y < fromXY.y {
		dir.y = -1
	}

	c := fromXY
	c.y += dir.y
	if dir.y < 0 {
		c.x += b.blockSize - 1
	}
	for dir.y*c.y < dir.y*toXY.y {
		b.elements[c.y][c.x].DrawVerticalLine()
		c.y += dir.y
	}
	b.elements[c.y][c.x].DrawJoint(dir.x > 0, dir.y < 0)

	stopBefore := 1
	if dir.x < 0 {
		stopBefore = b.blockSize
	}

	for dir.x*c.x < dir.x*toXY.x-stopBefore {
		c.x += dir.x
		b.elements[c.y][c.x].DrawHorizontalLine()
	}
	b.elements[c.y][c.x].PlaceArrow(dir.x < 0)
	return nil
}

func (b *Board) Render() string {
	rendered := ""
	for j := 0; j < b.h; j++ {
		for i := 0; i < b.w; i++ {
			rendered += b.elements[j][i].Render()
		}
		rendered += "\n"
	}
	return rendered
}
