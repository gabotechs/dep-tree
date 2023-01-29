package systems

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

const helpText = `
 ____         _ __       _
|  _ \   ___ |  _ \    _| |_  _ __  ___   ___ 
| | | | / _ \| |_) |  |_   _||  __|/ _ \ / _ \
| |_| ||  __/| .__/     | |  | |  |  __/|  __/
|____/  \__| |_|        | \__|_|   \___| \___|

Welcome to dep-tree's help section.

j      -> move one step down
k      -> move one step up
Ctrl d -> move half page down
Ctrl u -> move half page up
Enter  -> select the current node as the root node
q      -> navigate backwards on selected nodes or quit
h      -> show this help section
`

type ShouldQuit struct{}

func (s *ShouldQuit) Error() string {
	return "Should Quit"
}

func IsShouldQuit(err error) bool {
	_, ok := err.(*ShouldQuit)
	return ok
}

func navigate(s *State) error {
	if s.SelectedId == "" {
		return nil
	}
	err := s.Screen.Suspend()
	if err != nil {
		return err
	}
	err = s.OnNavigate(s)
	if err != nil {
		return err
	}
	err = s.Screen.Resume()
	if err != nil {
		return err
	}
	// NOTE: just to trigger an update.
	return s.Screen.PostEvent(&tcell.EventTime{})
}

func helpScreen(mockScreen tcell.Screen) error {
	var s tcell.Screen
	if mockScreen != nil {
		s = mockScreen
	} else {
		var err error
		s, err = tcell.NewScreen()
		if err != nil {
			return err
		}
	}
	err := s.Init()
	if err != nil {
		return err
	}
	lines := strings.Split(helpText, "\n")
	for y, line := range lines {
		for x := 0; x < len(line); x++ {
			s.SetContent(x, y, rune(line[x]), nil, defaultStyle)
		}
	}
	s.Show()
	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Rune() == 'q' {
				s.Fini()
				return nil
			}
		}
	}
}

func help(s *State) error {
	err := s.Screen.Suspend()
	if err != nil {
		return err
	}
	err = helpScreen(nil)
	if err != nil {
		return err
	}
	err = s.Screen.Resume()
	if err != nil {
		return err
	}
	// NOTE: just to trigger an update.
	return s.Screen.PostEvent(&tcell.EventTime{})
}

func RuntimeSystem(s *State) error {
	switch ev := s.Event.(type) {
	case *tcell.EventResize:
		s.Screen.Sync()
	case *tcell.EventInterrupt:
		s.Screen.Fini()
		return nil
	case *tcell.EventKey:
		switch {
		case ev.Rune() == 'q':
			s.Screen.Fini()
			return &ShouldQuit{}
		case ev.Rune() == 'h':
			return help(s)
		case ev.Key() == tcell.KeyEnter:
			return navigate(s)
		}
	}
	return nil
}
