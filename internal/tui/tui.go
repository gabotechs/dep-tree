package tui

import (
	"github.com/gabotechs/dep-tree/internal/tree"
	"github.com/gdamore/tcell/v2"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/ecs"
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
	parserBuilder dep_tree.NodeParserBuilder[T],
	screen tcell.Screen,
	isRootNavigation bool,
	tickChan chan bool,
) error {
	parser, err := parserBuilder(files)
	if err != nil {
		return err
	}
	dt := dep_tree.NewDepTree(parser, files).WithStdErrLoader()
	err = dt.LoadGraph()
	if err != nil {
		return err
	}
	dt.LoadCycles()
	t, err := tree.NewTree(dt)
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
			return Loop[T]([]string{s.SelectedId}, parserBuilder, screen, false, tickChan)
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
