package graphics

import (
	"errors"

	"github.com/gabotechs/dep-tree/internal/utils"
)

type LineTracer struct {
	slices []utils.Vector
	tags   map[string]string
}

func NewLineTracer(start utils.Vector) *LineTracer {
	return &LineTracer{
		slices: []utils.Vector{start},
		tags:   make(map[string]string),
	}
}

func (l *LineTracer) WithTags(tags map[string]string) *LineTracer {
	utils.Merge(l.tags, tags)
	return l
}

func (l *LineTracer) MoveVertical(reverse bool) utils.Vector {
	last := l.slices[len(l.slices)-1]
	current := last
	current.Y += utils.Bool2Int(!reverse)
	l.slices = append(l.slices, current)
	return current
}

func (l *LineTracer) MoveHorizontal(reverse bool) utils.Vector {
	last := l.slices[len(l.slices)-1]
	current := last
	current.X += utils.Bool2Int(!reverse)
	l.slices = append(l.slices, current)
	return current
}

func (l *LineTracer) Dump(matrix *Matrix) error {
	var lastCell *LinesCell
	for i := 1; i < len(l.slices); i++ {
		from := l.slices[i-1]
		to := l.slices[i]
		fromTo := to.Minus(from)
		if fromTo.X != 0 && fromTo.Y != 0 {
			return errors.New("cannot draw diagonal lines")
		}
		lines := LinesCell{
			l: fromTo.X < 0,
			r: fromTo.X > 0,
			t: fromTo.Y < 0,
			b: fromTo.Y > 0,
		}

		if lastCell == nil {
			lastCell = &lines
			startCellStack := matrix.Cell(from)
			if startCellStack != nil {
				startCellStack.add(NewTaggedCell(lastCell).WithTags(l.tags))
			}
		} else {
			if lines.l {
				lastCell.l = true
			}
			if lines.t {
				lastCell.t = true
			}
			if lines.r {
				lastCell.r = true
			}
			if lines.b {
				lastCell.b = true
			}
		}

		endCellStack := matrix.Cell(to)
		newCell := LinesCell{
			l: fromTo.X > 0,
			r: fromTo.X < 0,
			t: fromTo.Y > 0,
			b: fromTo.Y < 0,
		}

		if endCellStack != nil {
			endCellStack.add(NewTaggedCell(&newCell).WithTags(l.tags))
		}
		lastCell = &newCell
	}
	return nil
}
