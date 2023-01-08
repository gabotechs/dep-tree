package state

import (
	"github.com/gdamore/tcell/v2"
)

type State interface {
	Action(key tcell.Event)
	Update()
}

type States []State

func (s *States) Update() {
	for _, state := range *s {
		state.Update()
	}
}

func (s *States) Action(key tcell.Event) {
	for _, state := range *s {
		state.Action(key)
	}
}
