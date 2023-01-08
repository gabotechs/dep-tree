package tui

import (
	"github.com/gdamore/tcell/v2"

	"dep-tree/internal/utils"
)

var defaultStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
var primaryStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
var secondaryStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorBlue)

type Screen struct {
	tcell.Screen
}

func NewScreen() (*Screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	} else if err = s.Init(); err != nil {
		return nil, err
	}
	s.SetStyle(defaultStyle)

	return &Screen{Screen: s}, nil
}

func (s *Screen) Size() utils.Vector {
	w, h := s.Screen.Size()
	return utils.Vec(w, h)
}
