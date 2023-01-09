package systems

import (
	"github.com/gdamore/tcell/v2"

	"dep-tree/internal/utils"
)

const horizontalMargin = 0.5

type SpatialState struct {
	ScreenSize utils.Vector
	Offset     utils.Vector
	MaxY       int
}

func computeScreenOffset(s *State, ss *SpatialState) {
	verticalUpperLimit := ss.Offset.Y
	verticalLowerLimit := ss.Offset.Y + ss.ScreenSize.Y

	if s.Cursor.Y < verticalUpperLimit {
		ss.Offset.Y += s.Cursor.Y - verticalUpperLimit
	} else if s.Cursor.Y >= verticalLowerLimit {
		ss.Offset.Y += s.Cursor.Y - verticalLowerLimit + 1
	}

	horizontalUpperLimit := ss.Offset.X
	horizontalLowerLimit := ss.Offset.X + int(float32(ss.ScreenSize.X)*horizontalMargin)

	if s.Cursor.X < horizontalUpperLimit {
		ss.Offset.X += s.Cursor.X - horizontalUpperLimit
	} else if s.Cursor.X > horizontalLowerLimit {
		ss.Offset.X += s.Cursor.X - horizontalLowerLimit
	}
}

func up(n int, s *State, ss *SpatialState) {
	s.Cursor.Y -= n
	s.Cursor.Y = utils.Clamp(0, s.Cursor.Y, ss.MaxY)
}

func down(n int, s *State, ss *SpatialState) {
	s.Cursor.Y += n
	s.Cursor.Y = utils.Clamp(0, s.Cursor.Y, ss.MaxY)
}

func SpatialSystem(s *State, ss *SpatialState) {
	computeScreenOffset(s, ss)

	switch key := s.Event.(type) {
	case *tcell.EventResize:
		s.Screen.Sync()
		ss.ScreenSize.X, ss.ScreenSize.Y = key.Size()
	case *tcell.EventKey:
		switch key.Rune() {
		case 'j':
			down(1, s, ss)
		case 'k':
			up(1, s, ss)
		}

		switch key.Key() {
		case tcell.KeyDown:
			down(1, s, ss)
		case tcell.KeyUp:
			up(1, s, ss)
		case tcell.KeyCtrlU:
			up(ss.ScreenSize.Y/2, s, ss)
		case tcell.KeyCtrlD:
			down(ss.ScreenSize.Y/2, s, ss)
		}
	}
}
