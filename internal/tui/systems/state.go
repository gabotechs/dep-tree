package systems

import (
	"github.com/gdamore/tcell/v2"

	"github.com/gabotechs/dep-tree/internal/utils"
)

type State struct {
	Event            tcell.Event
	Screen           tcell.Screen
	SelectedId       string
	Cursor           utils.Vector
	OnNavigate       func(this *State) error
	IsRootNavigation bool
}
