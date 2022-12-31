package state

import (
	"github.com/gdamore/tcell/v2"

	"dep-tree/internal/board/graphics"
	"dep-tree/internal/graph"
	"dep-tree/internal/utils"
)

const horizontalMargin = 0.2
const verticalMargin = 0.3

type State struct {
	cells      [][]graphics.CellStack
	ScreenSize utils.Vector

	Cursor     utils.Vector
	Offset     utils.Vector
	SelectedId string
	ShouldQuit bool
}

func NewState(cells [][]graphics.CellStack, screenSize utils.Vector) *State {
	s := &State{cells: cells, ScreenSize: screenSize}
	s.update()
	return s
}

func (s *State) Quit() {
	s.ShouldQuit = true
}

func (s *State) update() {
	s.computeCursor()
	s.computeScreenOffset()
}

func (s *State) computeScreenOffset() {
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

func (s *State) computeCursor() {
	if s.Cursor.Y < len(s.cells) {
		for j := range s.cells[s.Cursor.Y] {
			if nodeId := s.cells[s.Cursor.Y][j].Tag(graph.NodeIdTag); nodeId != "" {
				s.Cursor.X = j
				s.SelectedId = nodeId
				return
			}
		}
	}
	s.SelectedId = ""
}

type RenderInfo struct {
	Position   utils.Vector
	Char       rune
	IsSelected bool
}

func (s *State) ForEachCell(f func(info RenderInfo)) {
	for i := s.Offset.Y; i < s.ScreenSize.Y+s.Offset.Y; i++ {
		if i >= len(s.cells) || i < 0 {
			break
		}
		for j := s.Offset.X; j < s.ScreenSize.X+s.Offset.X; j++ {
			if j >= len(s.cells[i]) || j < 0 {
				break
			}
			cell := s.cells[i][j]
			priorityTags := map[string]string{}
			isSelected := false
			if cell.Is(graph.NodeIdTag, s.SelectedId) {
				isSelected = true
			} else if cell.Is(graph.ConnectorOriginNodeIdTag, s.SelectedId) {
				isSelected = true
				priorityTags[graph.ConnectorOriginNodeIdTag] = s.SelectedId
			}
			f(RenderInfo{
				Position:   utils.Vec(j-s.Offset.X, i-s.Offset.Y),
				Char:       cell.Render(priorityTags),
				IsSelected: isSelected,
			})
		}
	}
}

func (s *State) SetSize(size utils.Vector) {
	s.ScreenSize = size
}

func (s *State) up(n int) {
	s.Cursor.Y -= n
	s.Cursor.Y = utils.Clamp(0, s.Cursor.Y, len(s.cells)-1)
}

func (s *State) down(n int) {
	s.Cursor.Y += n
	s.Cursor.Y = utils.Clamp(0, s.Cursor.Y, len(s.cells)-1)
}

func (s *State) Key(key *tcell.EventKey) {
	switch key.Rune() {
	case 'q':
		s.ShouldQuit = true
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
	s.update()
}

func (s *State) SetCursor(cursor utils.Vector) {
	s.Cursor = cursor
}

func (s *State) SetOffset(offset utils.Vector) {
	s.Offset = offset
}
