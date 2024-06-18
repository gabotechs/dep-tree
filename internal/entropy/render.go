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

//go:embed generated/index.html
var index []byte

const ToReplace = `"__INLINE_DATA",{}`
const ReplacePrefix = `"__INLINE_DATA",`

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
	// os.WriteFile("web/src/testData.json", marshaled, os.ModePerm)
	rendered := bytes.ReplaceAll(index, []byte(ToReplace), append([]byte(ReplacePrefix), marshaled...))
	var temp string
	if cfg.RenderPath != "" {
		temp = cfg.RenderPath
	} else {
		temp = filepath.Join(os.TempDir(), "index.html")
	}
	err = os.WriteFile(temp, rendered, 0o600)
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
