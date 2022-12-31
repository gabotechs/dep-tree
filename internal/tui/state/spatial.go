package state

import (
	"github.com/gdamore/tcell/v2"

	"dep-tree/internal/utils"
)

const horizontalMargin = 0.2
const verticalMargin = 0.3

type SpatialState struct {
	ScreenSize utils.Vector
	Cursor     *utils.Vector
	Offset     utils.Vector
	SelectedId string
	maxY       int
}

func NewSpatialState(
	cursor *utils.Vector,
	screenSize utils.Vector,
	maxY int,
) *SpatialState {
	s := &SpatialState{Cursor: cursor, ScreenSize: screenSize, maxY: maxY}
	return s
}

func (s *SpatialState) computeScreenOffset() {
	verticalUpperLimit := s.Offset.Y
	verticalLowerLimit := s.Offset.Y + int(float32(s.ScreenSize.Y)*verticalMargin)

	if s.Cursor.Y < verticalUpperLimit {
		s.Offset.Y += s.Cursor.Y - verticalUpperLimit
	} else if s.Cursor.Y > verticalLowerLimit {
		s.Offset.Y += s.Cursor.Y - verticalLowerLimit
	}

	horizontalUpperLimit := s.Offset.X
	horizontalLowerLimit := s.Offset.X + int(float32(s.ScreenSize.X)*horizontalMargin)

	if s.Cursor.X < horizontalUpperLimit {
		s.Offset.X += s.Cursor.X - horizontalUpperLimit
	} else if s.Cursor.X > horizontalLowerLimit {
		s.Offset.X += s.Cursor.X - horizontalLowerLimit
	}
}

func (s *SpatialState) SetSize(size utils.Vector) {
	s.ScreenSize = size
}

func (s *SpatialState) up(n int) {
	s.Cursor.Y -= n
	s.Cursor.Y = utils.Clamp(0, s.Cursor.Y, s.maxY)
}

func (s *SpatialState) down(n int) {
	s.Cursor.Y += n
	s.Cursor.Y = utils.Clamp(0, s.Cursor.Y, s.maxY)
}

func (s *SpatialState) SetOffset(offset utils.Vector) {
	s.Offset = offset
}

func (s *SpatialState) Action(ev tcell.Event) {
	switch key := ev.(type) {
	case *tcell.EventResize:
		s.ScreenSize.X, s.ScreenSize.Y = key.Size()
	case *tcell.EventKey:
		switch key.Rune() {
		case 'j':
			s.down(1)
		case 'k':
			s.up(1)
		}

		switch key.Key() {
		case tcell.KeyDown:
			s.down(1)
		case tcell.KeyUp:
			s.up(1)
		case tcell.KeyCtrlU:
			s.up(s.ScreenSize.Y)
		case tcell.KeyCtrlD:
			s.down(s.ScreenSize.Y)
		}
	}
}

func (s *SpatialState) Update() {
	s.computeScreenOffset()
}
