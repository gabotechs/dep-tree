package tui

import (
	"github.com/gdamore/tcell/v2"

	"github.com/gabotechs/dep-tree/internal/ecs"
	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/tree"
	"github.com/gabotechs/dep-tree/internal/tui/systems"
	"github.com/gabotechs/dep-tree/internal/utils"
)

var style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

func initScreen() (tcell.Screen, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	err = screen.Init()
	if err != nil {
		return nil, err
	}
	screen.SetStyle(style)
	return screen, nil
}

func Loop[T any](
	files []string,
	parser graph.NodeParser[T],
	screen tcell.Screen,
	isRootNavigation bool,
	tickChan chan bool,
	callbacks graph.LoadCallbacks[T],
) error {
	t, err := tree.NewTree(files, parser, callbacks)
	if err != nil {
		return err
	}
	board, err := t.Render()
	if err != nil {
		return err
	}
	cells, err := board.Cells()
	if err != nil {
		return err
	}
	if screen == nil {
		if screen, err = initScreen(); err != nil {
			return err
		}
	}

	renderState := &systems.RenderState{
		Cells:  cells,
		Errors: make(map[string][]error),
	}
	for _, n := range t.Nodes {
		renderState.Errors[n.Node.Id] = n.Node.Errors
	}
	spatialState := &systems.SpatialState{
		ScreenSize: utils.Vec(screen.Size()),
		Offset:     utils.Vec(0, 0),
		MaxY:       len(cells) - 1,
	}
	globalState := &systems.State{
		Cursor:           utils.Vec(0, 0),
		Screen:           screen,
		SelectedId:       "",
		Event:            nil,
		IsRootNavigation: isRootNavigation,
		OnNavigate: func(s *systems.State) error {
			return Loop[T]([]string{s.SelectedId}, parser, screen, false, tickChan, nil)
		},
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
		screen.Clear()
		err = world.Update()
		switch {
		case systems.IsShouldQuit(err):
			return nil
		case err != nil:
			return err
		}
		screen.Show()

		if tickChan != nil {
			tickChan <- true
		}
		globalState.Event = screen.PollEvent()
	}
}
