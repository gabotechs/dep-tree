package graphics

import (
	"fmt"

	"github.com/gabotechs/dep-tree/internal/utils"
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

func (m *Matrix) H() int {
	return m.h
}

func (m *Matrix) W() int {
	return m.w
}

func (m *Matrix) ExpandRight(n int) {
	for row := range m.elements {
		for i := 0; i < n; i++ {
			m.elements[row] = append(m.elements[row], CellStack{})
		}
	}
	m.w += n
}

func (m *Matrix) Cell(v utils.Vector) *CellStack {
	if v.Y >= 0 && v.X >= 0 && v.Y < len(m.elements) && v.X < len(m.elements[v.Y]) {
		return &m.elements[v.Y][v.X]
	} else {
		return nil
	}
}

func (m *Matrix) rayCast(
	origin utils.Vector,
	dir utils.Vector,
	query map[string]func(string) bool,
	length int,
) (bool, error) {
	for i := 0; i < length+1; i++ {
		cur := origin
		cur.X += dir.X * i
		cur.Y += dir.Y * i

		cellStack := m.Cell(cur)

		if cellStack == nil {
			if i == 0 {
				return false, fmt.Errorf("cannot ray cast in origin (%d, %d) because it is out of bounds", origin.X, origin.Y)
			} else {
				return false, nil
			}
		}

		for queryTag, queryFunction := range query {
			if value, ok := cellStack.Tags()[queryTag]; ok {
				if queryFunction(value) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (m *Matrix) RayCastVertical(
	origin utils.Vector,
	query map[string]func(string) bool,
	length int,
) (bool, error) {
	dir := 1
	if length < 0 {
		dir = -1
		length = -length
	}
	return m.rayCast(origin, utils.Vec(0, dir), query, length)
}

func (m *Matrix) Cells() [][]CellStack {
	return m.elements
}

func (m *Matrix) Render() string {
	rendered := ""
	for j := 0; j < m.h; j++ {
		for i := 0; i < m.w; i++ {
			rendered += string(m.elements[j][i].Render(nil))
		}
		rendered += "\n"
	}
	return rendered
}
