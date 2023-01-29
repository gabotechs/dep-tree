package systems

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/require"
)

func TestHelpScreen(t *testing.T) {
	a := require.New(t)
	mockScreen := tcell.NewSimulationScreen("")

	wait := make(chan error)

	go func() {
		wait <- helpScreen(mockScreen)
	}()

	time.Sleep(time.Millisecond * 100)
	mockScreen.InjectKey(tcell.Key(int16('q')), 'q', tcell.ModMask(0))

	err := <-wait
	a.NoError(err)
}
