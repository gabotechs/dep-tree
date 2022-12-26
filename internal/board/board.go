package board

import (
	"fmt"
	"strings"

	"github.com/elliotchance/orderedmap/v2"

	"dep-tree/internal/board/graphics"
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

func (b *Board) makeMatrix() (*graphics.Matrix, error) {
	matrix := graphics.NewMatrix(b.w, b.h)

	for _, k := range b.blocks.Keys() {
		block, _ := b.blocks.Get(k)
		err := block.Render(matrix)
		if err != nil {
			return matrix, fmt.Errorf("error rendering blockChar %s: %w", block.Label, err)
		}
	}

	for _, k := range b.connectors.Keys() {
		connector, _ := b.connectors.Get(k)
		err := connector.Render(matrix)
		if err != nil {
			return matrix, fmt.Errorf(
				"error rendering connector from %s to %s: %w",
				strings.TrimSpace(connector.from.Label),
				strings.TrimSpace(connector.to.Label),
				err,
			)
		}
	}
	return matrix, nil
}

func (b *Board) Cells() ([][]graphics.CellStack, error) {
	matrix, err := b.makeMatrix()
	return matrix.Cells(), err
}

func (b *Board) Render() (string, error) {
	matrix, err := b.makeMatrix()
	return matrix.Render(), err
}
