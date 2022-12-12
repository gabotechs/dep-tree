package board

import (
	"fmt"

	"github.com/elliotchance/orderedmap/v2"

	"dep-tree/internal/graphics"
)

type Board struct {
	w          int
	h          int
	blocks     *orderedmap.OrderedMap[string, *Block]
	connectors *orderedmap.OrderedMap[string, *Connector]
}

func MakeBoard() *Board {
	return &Board{
		blocks:     orderedmap.NewOrderedMap[string, *Block](),
		connectors: orderedmap.NewOrderedMap[string, *Connector](),
	}
}

func (b *Board) Render() (string, error) {
	// 1. Create Cell matrix.
	matrix := graphics.NewMatrix(b.w, b.h)

	// 2. Render blocks.
	for _, k := range b.blocks.Keys() {
		block, _ := b.blocks.Get(k)
		err := block.Render(matrix)
		if err != nil {
			return "", fmt.Errorf("error rendering block %s: %w", block.Id, err)
		}
	}

	// 3. Render connectors.
	for _, k := range b.connectors.Keys() {
		connector, _ := b.connectors.Get(k)
		err := connector.Render(matrix)
		if err != nil {
			return "", fmt.Errorf("error rendering connector from %s to %s: %w", connector.from.Id, connector.to.Id, err)
		}
	}

	// 4. dump Cells to a string.
	return matrix.Render(), nil
}
