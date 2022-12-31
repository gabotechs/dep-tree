package state

import (
	"github.com/gdamore/tcell/v2"
)

type State interface {
	Action(key tcell.Event)
	Update()
}
