package graphics

import (
	"errors"

	"dep-tree/internal/utils"
	"dep-tree/internal/vector"
)

type LineTracer struct {
	slices []vector.Vector
}

func NewLineTracer(start vector.Vector) *LineTracer {
	return &LineTracer{
		slices: []vector.Vector{start},
	}
}

func (l *LineTracer) MoveVertical(reverse bool) vector.Vector {
	last := l.slices[len(l.slices)-1]
	current := last
	current.Y += utils.Bool2Int(!reverse)
	l.slices = append(l.slices, current)
	return current
}

func (l *LineTracer) MoveHorizontal(reverse bool) vector.Vector {
	last := l.slices[len(l.slices)-1]
	current := last
	current.X += utils.Bool2Int(!reverse)
	l.slices = append(l.slices, current)
	return current
}

func (l *LineTracer) Dump(matrix *Matrix) error {
	var lastCell *Cell
	for i := 1; i < len(l.slices); i++ {
		from := l.slices[i-1]
		to := l.slices[i]
		fromTo := to.Minus(from)
		if fromTo.X != 0 && fromTo.Y != 0 {
			return errors.New("cannot draw diagonal lines")
		}
		lines := Lines{
			l: fromTo.X < 0,
			r: fromTo.X > 0,
			t: fromTo.Y < 0,
			b: fromTo.Y > 0,
		}

		if lastCell == nil {
			lastCell = LinesCell(lines)
			startCellStack := matrix.Cell(from)
			if startCellStack != nil {
				startCellStack.add(lastCell)
			}
		} else {
			lastCell.mergeLines(lines)
		}

		endCellStack := matrix.Cell(to)
		newCell := LinesCell(Lines{
			l: fromTo.X > 0,
			r: fromTo.X < 0,
			t: fromTo.Y > 0,
			b: fromTo.Y < 0,
		})

		if endCellStack != nil {
			endCellStack.add(newCell)
		}
		lastCell = newCell
	}
	return nil
}
