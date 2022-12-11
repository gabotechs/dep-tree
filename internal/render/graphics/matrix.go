package graphics

import (
	"fmt"

	"dep-tree/internal/vector"
)

type Matrix struct {
	w        int
	h        int
	elements [][]CellStack
}

func NewMatrix(w int, h int) *Matrix {
	elements := make([][]CellStack, h)
	for i := range elements {
		elements[i] = make([]CellStack, w)
	}
	return &Matrix{
		elements: elements,
		w:        w,
		h:        h,
	}
}

func (m *Matrix) Cell(v vector.Vector) *CellStack {
	return &m.elements[v.Y][v.X]
}

func (m *Matrix) rayCast(
	origin vector.Vector,
	dir vector.Vector,
	query map[string]func(string) bool,
	length int,
) (bool, error) {
	cur := origin
	for i := 0; i < length; i++ {
		cur.X += dir.X * i
		cur.Y += dir.Y * i

		if cur.X < 0 || cur.Y < 0 {
			if i == 0 {
				return false, fmt.Errorf("cannot ray cast in origin (%d, %d) because it is out of bounds", origin.X, origin.Y)
			} else {
				return false, nil
			}
		}

		if cur.X >= m.w || cur.Y >= m.h {
			if i == 0 {
				return false, fmt.Errorf("cannot ray cast in origin (%d, %d) because it is out of bounds", origin.X, origin.Y)
			} else {
				return false, nil
			}
		}

		cellStack := m.elements[cur.Y][cur.X]

		for queryTag, queryFunction := range query {
			if value, ok := cellStack.tags[queryTag]; ok {
				if queryFunction(value) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (m *Matrix) RayCastVertical(
	origin vector.Vector,
	query map[string]func(string) bool,
	length int,
) (bool, error) {
	dir := 1
	if length < 0 {
		dir = -1
	}
	return m.rayCast(origin, vector.Vec(0, dir), query, length)
}

func (m *Matrix) Render() string {
	rendered := ""
	for j := 0; j < m.h; j++ {
		for i := 0; i < m.w; i++ {
			rendered += m.elements[j][i].Render()
		}
		rendered += "\n"
	}
	return rendered
}
