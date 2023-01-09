package systems

import (
	"strings"

	"github.com/gdamore/tcell/v2"

	"dep-tree/internal/board/graphics"
	"dep-tree/internal/graph"
	"dep-tree/internal/utils"
)

type RenderState struct {
	Cells [][]graphics.CellStack
}

func computeCursor(s *State, rs *RenderState) {
	if s.Cursor.Y < len(rs.Cells) {
		for j := range rs.Cells[s.Cursor.Y] {
			if nodeId := rs.Cells[s.Cursor.Y][j].Tag(graph.NodeIdTag); nodeId != "" {
				s.Cursor.X = j
				s.SelectedId = nodeId
				return
			}
		}
	}
	s.SelectedId = ""
}

var defaultStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
var primaryStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
var secondaryStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorDarkCyan)

func forEachCell(
	s *State,
	rs *RenderState,
	ss *SpatialState,
) {
	from := utils.Vec(ss.Offset.X, ss.Offset.Y)
	to := utils.Vec(ss.Offset.X+ss.ScreenSize.X, ss.Offset.Y+ss.ScreenSize.Y)

	for i := from.Y; i < to.Y; i++ {
		if i >= len(rs.Cells) || i < 0 {
			break
		}
		for j := from.X; j < to.X; j++ {
			if j >= len(rs.Cells[i]) || j < 0 {
				break
			}
			cell := rs.Cells[i][j]
			priorityTags := map[string]string{}
			style := defaultStyle
			switch {
			case s.SelectedId == "":
				// nothing here.
			case cell.Is(graph.NodeIdTag, s.SelectedId):
				style = primaryStyle
			case cell.Is(graph.ConnectorOriginNodeIdTag, s.SelectedId):
				style = primaryStyle
				priorityTags[graph.ConnectorOriginNodeIdTag] = s.SelectedId
			case strings.Contains(cell.Tag(graph.NodeParentsTag), s.SelectedId):
				style = secondaryStyle
			}

			s.Screen.SetContent(
				j-from.X,
				i-from.Y,
				cell.Render(priorityTags),
				nil,
				style,
			)
		}
	}
}

func RenderSystem(s *State, rs *RenderState, ss *SpatialState) {
	computeCursor(s, rs)
	forEachCell(s, rs, ss)
}
