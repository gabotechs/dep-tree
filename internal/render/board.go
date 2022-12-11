package render

import (
	"github.com/elliotchance/orderedmap/v2"

	"dep-tree/internal/render/graphics"
	"dep-tree/internal/utils"
)

type BoardOptions struct {
	Indent int
}

type Board struct {
	w          int
	h          int
	options    BoardOptions
	blocks     *orderedmap.OrderedMap[string, *Block]
	connectors *orderedmap.OrderedMap[string, *Connector]
}

func MakeBoard(options BoardOptions) *Board {
	return &Board{
		options: BoardOptions{
			Indent: utils.Clamp(2, options.Indent, 8),
		},
		blocks:     orderedmap.NewOrderedMap[string, *Block](),
		connectors: orderedmap.NewOrderedMap[string, *Connector](),
	}
}

func (b *Board) Render() (string, error) {
	// 1. Create Cell matrix.
	elements := make([][]graphics.CellStack, b.h)
	for i := range elements {
		elements[i] = make([]graphics.CellStack, b.w)
	}

	// 2. Render blocks.
	for _, k := range b.blocks.Keys() {
		block, _ := b.blocks.Get(k)
		err := block.Render(elements)
		if err != nil {
			return "", err
		}
	}

	// 3. Render connectors.
	for _, k := range b.connectors.Keys() {
		connector, _ := b.connectors.Get(k)
		err := connector.Render(elements)
		if err != nil {
			return "", err
		}
	}

	// 4. dump Cells to a string.
	rendered := ""
	for j := 0; j < b.h; j++ {
		for i := 0; i < b.w; i++ {
			rendered += elements[j][i].Render()
		}
		rendered += "\n"
	}
	return rendered, nil
}
