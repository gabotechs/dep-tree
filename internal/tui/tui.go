package tui

import (
	"context"

	"github.com/gdamore/tcell/v2"

	"dep-tree/internal/dep_tree"
	"dep-tree/internal/ecs"
	"dep-tree/internal/tui/systems"
	"dep-tree/internal/utils"
)

var style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

type LoopTestOverrides struct {
	Screen tcell.SimulationScreen

	ManualPump bool
	Pump       func(event tcell.Event) error
}

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
	ctx context.Context,
	initial string,
	parserBuilder func(string) (dep_tree.NodeParser[T], error),
	testOverrides *LoopTestOverrides,
) error {
	parser, err := parserBuilder(initial)
	if err != nil {
		return err
	}
	ctx, dt, err := dep_tree.NewDepTree[T](ctx, parser)
	if err != nil {
		return err
	}
	board, err := dt.Render(parser.Display)
	if err != nil {
		return err
	}
	cells, err := board.Cells()
	if err != nil {
		return err
	}
	var screen tcell.Screen
	if testOverrides == nil || testOverrides.Screen == nil {
		if screen, err = initScreen(); err != nil {
			return err
		}
	} else {
		screen = testOverrides.Screen
	}

	renderState := &systems.RenderState{
		Cells:  cells,
		Errors: make(map[string][]error),
	}
	for _, n := range dt.Nodes {
		renderState.Errors[n.Node.Id] = n.Node.Errors
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
		OnNavigate: func(s *systems.State) error {
			return Loop[T](ctx, s.SelectedId, parserBuilder, testOverrides)
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
		if testOverrides != nil && testOverrides.ManualPump {
			testOverrides.Pump = makeTestPump(globalState, world)
			return nil
		} else {
			globalState.Event = screen.PollEvent()
		}
	}
}

func makeTestPump(s *systems.State, world *ecs.World) func(event tcell.Event) error {
	return func(event tcell.Event) error {
		s.Event = event
		s.Screen.Clear()
		err := world.Update()
		switch {
		case systems.IsShouldQuit(err):
			return nil
		case err != nil:
			return err
		}
		s.Screen.Show()
		return nil
	}
}
