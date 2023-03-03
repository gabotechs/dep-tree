package tui

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/require"

	"dep-tree/internal/js"
	"dep-tree/internal/language"
	"dep-tree/internal/rust"
	"dep-tree/internal/tui/systems"
	"dep-tree/internal/utils"
)

const tmp = "/tmp/dep-tree-tests"

const testPath = ".tui_test"

func TestTui(t *testing.T) {
	tests := []struct {
		Name       string
		Repo       string
		Tag        string
		Entrypoint string
		W          int
		H          int
		Keys       string
	}{
		{
			Name:       "react-stl-viewer",
			Repo:       "https://github.com/gabotechs/react-stl-viewer",
			Tag:        "2.2.4",
			Entrypoint: "src/index.ts",
			W:          40,
			H:          15,
		},
		{
			Name:       "react-gcode-viewer",
			Repo:       "https://github.com/gabotechs/react-gcode-viewer",
			Tag:        "2.2.4",
			Entrypoint: "src/GCodeViewer/GCodeModel.tsx",
			W:          40,
			H:          12,
		},
		{
			Name:       "graphql-js",
			Repo:       "https://github.com/graphql/graphql-js",
			Tag:        "v17.0.0-alpha.2",
			Entrypoint: "src/graphql.ts",
			W:          200,
			H:          130,
		},
		{
			Name:       "warp",
			Repo:       "https://github.com/seanmonstar/warp",
			Tag:        "v0.3.3",
			Entrypoint: "src/lib.rs",
			W:          100,
			H:          60,
		},
		{
			Name:       "quits",
			Repo:       "https://github.com/seanmonstar/warp",
			Tag:        "v0.3.3",
			Entrypoint: "src/lib.rs",
			W:          40,
			H:          30,
			Keys:       "q",
		},
	}

	_ = os.MkdirAll(testPath, os.ModePerm)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			repoPath := path.Join(tmp, path.Base(tt.Repo))
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
			err := screen.Init()
			a.NoError(err)
			screen.SetSize(tt.W, tt.H)
			testOverrides := &LoopTestOverrides{
				Screen:     screen,
				ManualPump: true,
			}
			if utils.EndsWith(entrypointPath, js.Extensions) {
				err := Loop[js.Data](
					context.Background(),
					entrypointPath,
					language.ParserBuilder(js.MakeJsLanguage),
					testOverrides,
				)
				a.NoError(err)
			} else if utils.EndsWith(entrypointPath, rust.Extensions) {
				err := Loop[rust.Data](
					context.Background(),
					entrypointPath,
					language.ParserBuilder(rust.MakeRustLanguage),
					testOverrides,
				)
				a.NoError(err)
			}
			if tt.Keys != "" {
				for _, key := range strings.Split(tt.Keys, " ") {
					err := testOverrides.Pump(tCellKey(key))
					a.NoError(err)
				}
			}

			result := systems.PrintScreen(screen)
			utils.GoldenTest(
				t,
				path.Join(".tui_test", tt.Name+".txt"),
				result,
			)
		})
	}
}

func tCellKey(str string) tcell.Event {
	return tcell.NewEventKey(tcell.Key(str[0]), rune(str[0]), tcell.ModMask(0))
}
