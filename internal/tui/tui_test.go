package tui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/require"

	"dep-tree/internal/js"
	"dep-tree/internal/language"
)

const tmp = "/tmp/dep-tree-tests"

const testPath = ".tui_test"

func printScreen(s tcell.SimulationScreen) string {
	result := ""
	w, h := s.Size()
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			char, _, _, _ := s.GetContent(j, i)
			result += string(char)
		}
		result += "\n"
	}
	return result
}

func TestTui(t *testing.T) {
	tests := []struct {
		Name       string
		Repo       string
		Tag        string
		Entrypoint string
	}{
		{
			Name:       "react-stl-viewer",
			Repo:       "https://github.com/gabotechs/react-stl-viewer",
			Tag:        "2.2.4",
			Entrypoint: "src/index.ts",
		},
		{
			Name:       "react-gcode-viewer",
			Repo:       "https://github.com/gabotechs/react-gcode-viewer",
			Tag:        "2.2.4",
			Entrypoint: "src/GCodeViewer/GCodeModel.tsx",
		},
		{
			Name:       "graphql-js",
			Repo:       "https://github.com/graphql/graphql-js",
			Tag:        "v17.0.0-alpha.2",
			Entrypoint: "src/graphql.ts",
		},
	}

	_ = os.MkdirAll(testPath, os.ModePerm)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			repoPath := path.Join(tmp, path.Base(tt.Name))
			entrypointPath := path.Join(repoPath, tt.Entrypoint)
			if _, err := os.Stat(entrypointPath); err != nil {
				_ = os.RemoveAll(repoPath)
				_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
					URL:           tt.Repo,
					ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", tt.Tag)),
					SingleBranch:  true,
					Depth:         1,
					Progress:      os.Stdout,
				})
				a.NoError(err)
			}

			screen := tcell.NewSimulationScreen("")

			wait := make(chan error)

			go func() {
				wait <- Loop[js.Data](
					context.Background(),
					entrypointPath,
					language.ParserBuilder(js.MakeJsLanguage),
					screen,
				)
			}()
			ticker := time.NewTicker(time.Millisecond * 300)
			elapsed := 0

			result := ""

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
					result = printScreen(screen)
					screen.InjectKey(tcell.Key(int16('q')), 'q', tcell.ModMask(0))
					elapsed++
				}
			}

			a.NoError(err)

			expectedPath := path.Join(testPath, tt.Name+".txt")
			if _, err := os.Stat(expectedPath); err == nil {
				expected, err := os.ReadFile(expectedPath)
				a.NoError(err)
				a.Equal(string(expected), result)
			} else {
				err := os.WriteFile(expectedPath, []byte(result), os.ModePerm)
				a.NoError(err)
			}
		})
	}
}
