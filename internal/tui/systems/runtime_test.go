package systems

import (
	"path"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/require"

	"dep-tree/internal/utils"
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

	result := PrintScreen(mockScreen)
	utils.GoldenTest(t, path.Join(".runtime_system_test", "help.txt"), result)

	err := <-wait
	a.NoError(err)
}
