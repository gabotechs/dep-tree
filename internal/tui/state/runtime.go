package state

import "github.com/gdamore/tcell/v2"

type RunTimeState struct {
	ShouldQuit bool
}

func NewRuntimeState() *RunTimeState {
	return &RunTimeState{
		ShouldQuit: false,
	}
}

func (rts *RunTimeState) Quit() {
	rts.ShouldQuit = true
}

func (rts *RunTimeState) Action(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventInterrupt:
		rts.ShouldQuit = true
	case *tcell.EventKey:
		if ev.Rune() == 'q' {
			rts.ShouldQuit = true
		}
	}
}

func (rts *RunTimeState) Update() {

}
