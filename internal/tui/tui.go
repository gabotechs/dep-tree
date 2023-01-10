package tui

import (
	"context"

	"github.com/gdamore/tcell/v2"

	"dep-tree/internal/ecs"
	"dep-tree/internal/graph"
	"dep-tree/internal/tui/systems"
	"dep-tree/internal/utils"
)

var style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

func Loop[T any](
	ctx context.Context,
	entrypoint string,
	parser graph.NodeParser[T],
	screen tcell.Screen,
) error {
	if screen == nil {
		var err error
		screen, err = tcell.NewScreen()
		if err != nil {
			return err
		}
	}
	err := screen.Init()
	if err != nil {
		return err
	}
	screen.SetStyle(style)
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

	renderState := &systems.RenderState{
		Cells: cells,
	}
	spatialState := &systems.SpatialState{
		ScreenSize: utils.Vec(screen.Size()),
		Offset:     utils.Vec(0, 0),
		MaxY:       len(cells) - 1,
	}
	globalState := &systems.State{
		Cursor:     utils.Vec(0, 0),
		Screen:     screen,
		SelectedId: "",
		Event:      nil,
	}

	world := ecs.NewWorld().
		WithEntity(ecs.NewEntity().
			With(globalState).
			With(renderState).
			With(spatialState),
		).
		WithSystem(systems.SpatialSystem).
		WithSystem(systems.RenderSystem).
		WithSystem(systems.RuntimeSystem)

	for {
		err = world.Update()
		switch {
		case systems.IsShouldQuit(err):
			return nil
		case systems.IsShouldNavigate(err):
			_ = screen.Suspend()
			err = Loop[T](ctx, globalState.SelectedId, parser, nil)
			if err != nil {
				return err
			}
			_ = screen.Resume()
			_ = world.Update()
		case err != nil:
			return err
		}
		screen.Show()

		ev := screen.PollEvent()
		if _, ok := ev.(*tcell.EventResize); ok {
			screen.Sync()
		}
		globalState.Event = ev
	}
}
