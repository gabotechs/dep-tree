package tui

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/require"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/js"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python"
	"github.com/gabotechs/dep-tree/internal/rust"
	"github.com/gabotechs/dep-tree/internal/tui/systems"
	"github.com/gabotechs/dep-tree/internal/utils"
)

const tmp = "/tmp/dep-tree-tests"

const testPath = ".tui_test"

//nolint:gocyclo
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
			Name:       "jumps 4 downs",
			Repo:       "https://github.com/gabotechs/react-stl-viewer",
			Tag:        "2.2.4",
			Entrypoint: "src/index.ts",
			W:          40,
			H:          15,
			Keys:       "j j j j",
		},
		{
			Name:       "jumps 4 ups",
			Repo:       "https://github.com/gabotechs/react-stl-viewer",
			Tag:        "2.2.4",
			Entrypoint: "src/index.ts",
			W:          40,
			H:          15,
			Keys:       "k k k k",
		},
		{
			Name:       "jumps 5 down and navigates",
			Repo:       "https://github.com/gabotechs/react-stl-viewer",
			Tag:        "2.2.4",
			Entrypoint: "src/index.ts",
			W:          40,
			H:          15,
			Keys:       "j j j j j enter",
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
			Name:       "graphql-js with ctrl displacements",
			Repo:       "https://github.com/graphql/graphql-js",
			Tag:        "v17.0.0-alpha.2",
			Entrypoint: "src/graphql.ts",
			W:          100,
			H:          20,
			Keys:       "down down up ctrl-d ctrl-d ctrl-u",
		},
		{
			Name:       "warp",
			Repo:       "https://github.com/seanmonstar/warp",
			Tag:        "v0.3.3",
			Entrypoint: "src/lib.rs",
			W:          100,
			H:          60,
		},
	}

	_ = os.MkdirAll(testPath, os.ModePerm)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			repoPath := filepath.Join(tmp, path.Base(tt.Repo))
			entrypointPath := filepath.Join(repoPath, tt.Entrypoint)
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

			update := make(chan bool)
			finish := make(chan error)

			go func() {
				var parserBuilder dep_tree.NodeParserBuilder[language.FileInfo]
				switch {
				case utils.EndsWith(entrypointPath, js.Extensions):
					parserBuilder = language.ParserBuilder(js.MakeJsLanguage, nil, nil)
				case utils.EndsWith(entrypointPath, rust.Extensions):
					parserBuilder = language.ParserBuilder(rust.MakeRustLanguage, nil, nil)
				case utils.EndsWith(entrypointPath, python.Extensions):
					parserBuilder = language.ParserBuilder(python.MakePythonLanguage, nil, nil)
				}

				finish <- Loop[language.FileInfo](
					entrypointPath,
					parserBuilder,
					screen,
					true,
					update,
				)
			}()

			select {
			case <-update:
			case err = <-finish:
			}
			a.NoError(err)

			nQs := 1
			if tt.Keys != "" {
				for _, key := range strings.Split(tt.Keys, " ") {
					var e tcell.Event
					switch key {
					case "enter":
						nQs++
						e = tCellEnter()
					case "down":
						e = tCellDown()
					case "up":
						e = tCellUp()
					case "ctrl-u":
						e = tCellCtrlU()
					case "ctrl-d":
						e = tCellCtrlD()
					default:
						e = tCellKey(key)
					}
					err := screen.PostEvent(e)
					a.NoError(err)
					<-update
				}
			}

			result := systems.PrintScreen(screen)
			utils.GoldenTest(
				t,
				filepath.Join(".tui_test", tt.Name+".txt"),
				result,
			)
			for i := 0; i < nQs; i++ {
				err := screen.PostEvent(tCellKey("q"))
				a.NoError(err)
				select {
				case <-update:
				case err := <-finish:
					a.NoError(err)
				}
			}
		})
	}
}

func tCellDown() tcell.Event {
	return tcell.NewEventKey(tcell.KeyDown, ' ', tcell.ModNone)
}

func tCellUp() tcell.Event {
	return tcell.NewEventKey(tcell.KeyUp, ' ', tcell.ModNone)
}

func tCellCtrlU() tcell.Event {
	return tcell.NewEventKey(tcell.KeyCtrlU, ' ', tcell.ModNone)
}

func tCellCtrlD() tcell.Event {
	return tcell.NewEventKey(tcell.KeyCtrlD, ' ', tcell.ModNone)
}

func tCellKey(str string) tcell.Event {
	return tcell.NewEventKey(tcell.Key(str[0]), rune(str[0]), tcell.ModNone)
}

func tCellEnter() tcell.Event {
	return tcell.NewEventKey(tcell.KeyEnter, ' ', tcell.ModNone)
}
