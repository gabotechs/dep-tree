package entropy

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/language"
)

//go:embed index.html
var index []byte

const ToReplace = "const DATA = {}"
const ReplacePrefix = "const DATA = "

type RenderConfig struct {
	NoOpen    bool
	EnableGui bool
}

func Render(parser language.NodeParser, cfg RenderConfig) error {
	dt := dep_tree.NewDepTree(parser)
	err := dt.LoadGraph()
	if err != nil {
		return err
	}

	dt.LoadCycles()
	graph := makeGraph(dt, parser)
	graph.EnableGui = cfg.EnableGui
	marshaled, err := json.Marshal(graph)
	if err != nil {
		return err
	}
	rendered := bytes.ReplaceAll(index, []byte(ToReplace), append([]byte(ReplacePrefix), marshaled...))
	temp := path.Join(os.TempDir(), "index.html")
	err = os.WriteFile(temp, rendered, os.ModePerm)
	if err != nil {
		return err
	}
	if cfg.NoOpen {
		fmt.Println(temp)
		return nil
	} else {
		return openInBrowser(temp)
	}
}
