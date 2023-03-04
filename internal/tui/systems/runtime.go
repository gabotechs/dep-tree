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
	// NOTE: just to trigger an update.
	defer func() { _ = s.Screen.PostEvent(&tcell.EventTime{}) }()
	return s.OnNavigate(s)
}

func helpScreen(screen tcell.Screen) error {
	lines := strings.Split(helpText, "\n")
	screen.Clear()
	for y, line := range lines {
		for x := 0; x < len(line); x++ {
			screen.SetContent(x, y, rune(line[x]), nil, defaultStyle)
		}
	}
	screen.Show()
	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if ev.Rune() == 'q' {
				return nil
			}
		}
	}
}

func help(screen tcell.Screen) error {
	defer func() { _ = screen.PostEvent(&tcell.EventTime{}) }()
	return helpScreen(screen)
}

func RuntimeSystem(s *State) error {
	switch ev := s.Event.(type) {
	case *tcell.EventResize:
		s.Screen.Sync()
	case *tcell.EventInterrupt:
		if s.IsRootNavigation {
			s.Screen.Fini()
		}
		return &ShouldQuit{}
	case *tcell.EventKey:
		switch {
		case ev.Rune() == 'q':
			if s.IsRootNavigation {
				s.Screen.Fini()
			}
			return &ShouldQuit{}
		case ev.Rune() == 'h':
			return help(s.Screen)
		case ev.Key() == tcell.KeyEnter:
			return navigate(s)
		}
	}
	return nil
}
