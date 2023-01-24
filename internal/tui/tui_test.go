package tui

import (
	"context"
	"errors"
	"os"
	"path"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/require"

	"dep-tree/internal/dep_tree"
	"dep-tree/internal/js"
)

const tmp = "/tmp/dep-tree-tests"

func TestTui(t *testing.T) {
	tests := []struct {
		Name       string
		Repo       string
		Entrypoint string
	}{
		{
			Name:       "react-stl-viewer",
			Repo:       "https://github.com/gabotechs/react-stl-viewer",
			Entrypoint: "src/index.ts",
		},
		{
			Name:       "react-gcode-viewer",
			Repo:       "https://github.com/gabotechs/react-gcode-viewer",
			Entrypoint: "src/GCodeViewer/GCodeModel.tsx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			repoPath := path.Join(tmp, path.Base(tt.Name))
			entrypointPath := path.Join(repoPath, tt.Entrypoint)
			if _, err := os.Stat(entrypointPath); err != nil {
				_ = os.RemoveAll(repoPath)
				_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
					URL:      tt.Repo,
					Progress: os.Stdout,
				})
				a.NoError(err)
			}

			screen := tcell.NewSimulationScreen("")

			wait := make(chan error)

			go func() {
				wait <- Loop[js.Data](
					context.Background(),
					entrypointPath,
					func(s string) (dep_tree.NodeParser[js.Data], error) {
						return js.MakeJsParser(s)
					},
					screen,
				)
			}()
			ticker := time.NewTicker(time.Millisecond * 100)
			elapsed := 0

			var err error
		forLoop:
			for {
				select {
				case result := <-wait:
					err = result
					break forLoop
				case <-ticker.C:
					if elapsed > 100 {
						err = errors.New("timeout")
						break forLoop
					}
					screen.InjectKey(tcell.Key(int16('q')), 'q', tcell.ModMask(0))
					elapsed++
				}
			}

			a.NoError(err)
		})
	}
}
