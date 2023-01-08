package tui

import (
	"github.com/gdamore/tcell/v2"

	"dep-tree/internal/board"
	s "dep-tree/internal/tui/state"
	"dep-tree/internal/utils"
)

func Loop(b *board.Board) error {
	screen, err := NewScreen()
	if err != nil {
		return err
	}
	cells, err := b.Cells()
	if err != nil {
		return err
	}
	cursor := utils.Vec(0, 0)
	renderState := s.NewRenderState(&cursor, cells)
	runtimeState := s.NewRuntimeState()
	spatialState := s.NewSpatialState(&cursor, screen.Size(), len(cells)-1)

	renderState.Update()
	spatialState.Update()
	runtimeState.Update()

	for {
		screen.Clear()
		if runtimeState.ShouldQuit {
			screen.Fini()
			return nil
		}

		renderState.ForEachCell(
			utils.Vec(spatialState.Offset.X, spatialState.Offset.Y),
			utils.Vec(spatialState.Offset.X+spatialState.ScreenSize.X, spatialState.Offset.Y+spatialState.ScreenSize.Y),
			func(info s.RenderInfo) {
				style := defaultStyle
				if info.IsSelected {
					style = primaryStyle
				} else if info.IsHighlighted {
					style = secondaryStyle
				}

				screen.SetContent(
					info.Position.X,
					info.Position.Y,
					info.Char,
					nil,
					style,
				)
			})

		screen.Show()

		ev := screen.PollEvent()
		if _, ok := ev.(*tcell.EventResize); ok {
			screen.Sync()
		}

		renderState.Action(ev)
		spatialState.Action(ev)
		runtimeState.Action(ev)

		renderState.Update()
		spatialState.Update()
		runtimeState.Update()
	}
}
