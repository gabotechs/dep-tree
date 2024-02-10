package entropy

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/language"
)

//go:embed index.html
var index []byte

const ToReplace = "const DATA = {}"
const ReplacePrefix = "const DATA = "

type RenderConfig struct {
	NoOpen        bool
	EnableGui     bool
	RenderPath    string
	LoadCallbacks graph.LoadCallbacks[*language.FileInfo]
}

func Render(files []string, parser graph.NodeParser[*language.FileInfo], cfg RenderConfig) error {
	graph3d, err := makeGraph(files, parser, cfg.LoadCallbacks)
	if err != nil {
		return err
	}
	graph3d.EnableGui = cfg.EnableGui
	marshaled, err := json.Marshal(graph3d)
	if err != nil {
		return err
	}
	rendered := bytes.ReplaceAll(index, []byte(ToReplace), append([]byte(ReplacePrefix), marshaled...))
	var temp string
	if cfg.RenderPath != "" {
		temp = cfg.RenderPath
	} else {
		temp = filepath.Join(os.TempDir(), "index.html")
	}
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
