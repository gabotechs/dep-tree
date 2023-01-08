package tui

import (
	"context"

	"github.com/gdamore/tcell/v2"

	"dep-tree/internal/graph"
	s "dep-tree/internal/tui/state"
	"dep-tree/internal/utils"
)

func Loop[T any](ctx context.Context, entrypoint string, parser graph.NodeParser[T]) error {
	screen, err := NewScreen()
	if err != nil {
		return err
	}

	_, board, err := graph.RenderGraph(ctx, entrypoint, parser)
	if err != nil {
		return err
	}

	cells, err := board.Cells()
	if err != nil {
		return err
	}
	cursor := utils.Vec(0, 0)

	renderState := s.NewRenderState(&cursor, cells)
	runtimeState := s.NewRuntimeState()
	spatialState := s.NewSpatialState(&cursor, screen.Size(), len(cells)-1)

	states := s.States{
		renderState,
		spatialState,
		runtimeState,
	}

	states.Update()

	for {
		if runtimeState.ShouldQuit {
			screen.Fini()
			return nil
		} else if runtimeState.Next {
			if renderState.SelectedId != "" {
				err = screen.Suspend()
				if err != nil {
					return err
				}
				err = Loop[T](ctx, renderState.SelectedId, parser)
				if err != nil {
					return err
				}
				err = screen.Resume()
				if err != nil {
					return err
				}
			}
			runtimeState.Next = false
		}

		screen.Clear()

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

		states.Action(ev)
		states.Update()
	}
}
