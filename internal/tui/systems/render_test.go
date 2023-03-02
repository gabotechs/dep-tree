package systems

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/require"

	"dep-tree/internal/utils"
)

const renderTestFolder = ".render_system_test"

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
			Name: "Very long word",
			Errors: []error{
				errors.New("/Users/gabriel/GolandProjects/dep-tree/internal/tui/systems/render_test.go"),
			},
			ScreenSize: utils.Vec(100, 10),
		},
	}

	_ = os.MkdirAll(renderTestFolder, os.ModePerm)

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

			gather := make([]byte, 0)

			for y := 0; y < tt.ScreenSize.Y; y++ {
				for x := 0; x < tt.ScreenSize.X; x++ {
					c, _, _, _ := screen.GetContent(x, y)
					gather = append(gather, byte(c))
				}
				gather = append(gather, byte('\n'))
			}
			resultFile := path.Join(renderTestFolder, tt.Name+".txt")
			if utils.FileExists(resultFile) {
				expected, err := os.ReadFile(resultFile)
				a.NoError(err)
				a.Equal(string(expected), string(gather))
			} else {
				err := os.WriteFile(resultFile, gather, os.ModePerm)
				a.NoError(err)
			}
		})
	}
}
