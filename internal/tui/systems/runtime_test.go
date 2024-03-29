package systems

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/utils"
)

func TestHelpScreen(t *testing.T) {
	a := require.New(t)
	mockScreen := tcell.NewSimulationScreen("")
	err := mockScreen.Init()
	a.NoError(err)

	wait := make(chan error)

	go func() {
		wait <- helpScreen(mockScreen)
	}()

	time.Sleep(time.Millisecond * 100)
	result := PrintScreen(mockScreen)

	mockScreen.InjectKey(tcell.Key(int16('q')), 'q', tcell.ModMask(0))

	utils.GoldenTest(t, filepath.Join(".runtime_system_test", "help.txt"), result)

	err = <-wait
	a.NoError(err)
}
