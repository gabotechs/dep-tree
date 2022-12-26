package tui

import (
	"strconv"

	"github.com/gdamore/tcell/v2"

	"dep-tree/internal/board"
)

var defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
var selectedStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)

type State struct {
	selected int
}

func Loop(b *board.Board) error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	} else if err = s.Init(); err != nil {
		return err
	}
	s.SetStyle(defStyle)

	state := State{
		selected: 0,
	}

	for {
		s.Clear()
		cells, err := b.Cells()
		if err != nil {
			return err
		}
		for i := range cells {
			for j := range cells[i] {
				style := defStyle
				if cells[i][j].Is("nodeIndex", strconv.Itoa(state.selected)) {
					style = selectedStyle
				}
				s.SetContent(
					j,
					i,
					cells[i][j].Render(map[string]string{
						"nodeIndex": strconv.Itoa(state.selected),
					}),
					nil,
					style,
				)
			}
		}
		s.Show()

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			switch ev.Rune() {
			case 'q':
				return nil
			case 'j':
				state.selected++
			case 'k':
				if state.selected > 0 {
					state.selected--
				}
			}
		}
	}
}
