package state

import (
	"github.com/gdamore/tcell/v2"

	"dep-tree/internal/board/graphics"
	"dep-tree/internal/graph"
	"dep-tree/internal/utils"
)

type RenderState struct {
	cells      [][]graphics.CellStack
	SelectedId string
	Cursor     *utils.Vector
}

func NewRenderState(cursor *utils.Vector, cells [][]graphics.CellStack) *RenderState {
	return &RenderState{
		Cursor:     cursor,
		cells:      cells,
		SelectedId: "",
	}
}

type RenderInfo struct {
	Position   utils.Vector
	Char       rune
	IsSelected bool
}

func (rs *RenderState) computeCursor() {
	if rs.Cursor.Y < len(rs.cells) {
		for j := range rs.cells[rs.Cursor.Y] {
			if nodeId := rs.cells[rs.Cursor.Y][j].Tag(graph.NodeIdTag); nodeId != "" {
				rs.Cursor.X = j
				rs.SelectedId = nodeId
				return
			}
		}
	}
	rs.SelectedId = ""
}

func (rs *RenderState) ForEachCell(
	from utils.Vector,
	to utils.Vector,
	f func(info RenderInfo),
) {
	for i := from.Y; i < to.Y; i++ {
		if i >= len(rs.cells) || i < 0 {
			break
		}
		for j := from.X; j < to.X; j++ {
			if j >= len(rs.cells[i]) || j < 0 {
				break
			}
			cell := rs.cells[i][j]
			priorityTags := map[string]string{}
			isSelected := false
			if cell.Is(graph.NodeIdTag, rs.SelectedId) {
				isSelected = true
			} else if cell.Is(graph.ConnectorOriginNodeIdTag, rs.SelectedId) {
				isSelected = true
				priorityTags[graph.ConnectorOriginNodeIdTag] = rs.SelectedId
			}
			f(RenderInfo{
				Position:   utils.Vec(j-from.X, i-from.Y),
				Char:       cell.Render(priorityTags),
				IsSelected: isSelected,
			})
		}
	}
}

func (rs *RenderState) Action(key tcell.Event) {
	// nothing.
}

func (rs *RenderState) Update() {
	rs.computeCursor()
}
