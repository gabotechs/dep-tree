package entropy

import (
	"bytes"
	"context"
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

func Render(ctx context.Context, parser language.NodeParser, cfg RenderConfig) (context.Context, error) {
	dt := dep_tree.NewDepTree(parser)
	ctx, err := dt.LoadGraph(ctx)
	if err != nil {
		return ctx, err
	}

	dt.LoadCycles()
	graph := makeGraph(dt, parser)
	graph.EnableGui = cfg.EnableGui
	marshaled, err := json.Marshal(graph)
	if err != nil {
		return ctx, err
	}
	rendered := bytes.ReplaceAll(index, []byte(ToReplace), append([]byte(ReplacePrefix), marshaled...))
	temp := path.Join(os.TempDir(), "index.html")
	err = os.WriteFile(temp, rendered, os.ModePerm)
	if err != nil {
		return ctx, err
	}
	if cfg.NoOpen {
		fmt.Println(temp)
		return ctx, nil
	} else {
		return ctx, openInBrowser(temp)
	}
}
