package tui

import (
	"github.com/gdamore/tcell/v2"

	"dep-tree/internal/board"
	s "dep-tree/internal/tui/state"
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
	state := s.NewState(cells, screen.Size())

	for {
		screen.Clear()
		if state.ShouldQuit {
			screen.Fini()
			return nil
		}

		state.ForEachCell(func(info s.RenderInfo) {
			style := defStyle
			if info.IsSelected {
				style = selectedStyle
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

		switch ev := ev.(type) {
		case *tcell.EventInterrupt:
			state.Quit()
		case *tcell.EventResize:
			screen.Sync()
			state.SetSize(screen.Size())
		case *tcell.EventKey:
			state.Key(ev)
		}
	}
}
