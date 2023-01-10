package systems

import "github.com/gdamore/tcell/v2"

type ShouldQuit struct{}

func (s *ShouldQuit) Error() string {
	return "Should Quit"
}

func IsShouldQuit(err error) bool {
	_, ok := err.(*ShouldQuit)
	return ok
}

type ShouldNavigate struct{}

func (s *ShouldNavigate) Error() string {
	return "Should Navigate"
}

func IsShouldNavigate(err error) bool {
	_, ok := err.(*ShouldNavigate)
	return ok
}

func RuntimeSystem(s *State) error {
	switch ev := s.Event.(type) {
	case *tcell.EventInterrupt:
		s.Screen.Fini()
		return nil
	case *tcell.EventKey:
		if ev.Rune() == 'q' {
			s.Screen.Fini()
			return &ShouldQuit{}
		} else if ev.Key() == tcell.KeyEnter {
			if s.SelectedId == "" {
				return nil
			}
			return &ShouldNavigate{}
		}
	}
	return nil
}
