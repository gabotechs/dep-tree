package systems

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/utils"
)

func TestRenderSystem(t *testing.T) {
	tests := []struct {
		Name       string
		Errors     []error
		ScreenSize utils.Vector
	}{
		{
			Name: "Short error",
			Errors: []error{
				errors.New("this is an error"),
			},
			ScreenSize: utils.Vec(60, 3),
		},
		{
			Name: "Long error",
			Errors: []error{
				errors.New(`this is a very long error that probably does not fit in just one line, so it will be split in multiple lines`),
			},
			ScreenSize: utils.Vec(60, 10),
		},
		{
			Name: "Two errors",
			Errors: []error{
				errors.New("one error that is longer than the one below, but not that much"),
				errors.New("another short error"),
			},
			ScreenSize: utils.Vec(60, 10),
		},
		{
			Name: "Two equal errors",
			Errors: []error{
				errors.New("one error that is longer than the one below, but not that much"),
				errors.New("one error that is longer than the one below, but not that much"),
			},
			ScreenSize: utils.Vec(60, 10),
		},
		{
			Name: "Very long word",
			Errors: []error{
				errors.New("/Users/gabriel/GolandProjects/dep-tree/internal/tui/systems/render_test.go"),
			},
			ScreenSize: utils.Vec(100, 10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			screen := tcell.NewSimulationScreen("")
			err := screen.Init()
			a.NoError(err)
			screen.SetStyle(defaultStyle)
			screen.SetSize(tt.ScreenSize.X, tt.ScreenSize.Y)

			s := &State{
				SelectedId: "selected",
				Screen:     screen,
			}

			rs := &RenderState{
				Errors: map[string][]error{
					"selected": tt.Errors,
				},
			}

			ss := &SpatialState{
				ScreenSize: tt.ScreenSize,
			}

			renderError(s, rs, ss)

			gather := PrintScreen(screen)

			utils.GoldenTest(
				t,
				filepath.Join(".render_system_test", tt.Name+".txt"),
				gather,
			)
		})
	}
}
