package entropy

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/utils"
)

//go:embed index.html
var index string

const ToReplace = "const GRAPH = {}"
const ReplacePrefix = "const GRAPH = "

func Render(ctx context.Context, parser language.NodeParser) (context.Context, error) {
	dt := dep_tree.NewDepTree(parser)
	ctx, err := dt.LoadGraph(ctx)
	if err != nil {
		return ctx, err
	}
	marshaled, err := marshalGraph(dt.Graph, parser)
	if err != nil {
		return ctx, err
	}
	rendered := strings.ReplaceAll(index, ToReplace, ReplacePrefix+string(marshaled))
	temp := path.Join(os.TempDir(), "index.html")
	err = os.WriteFile(temp, []byte(rendered), os.ModePerm)
	if err != nil {
		return ctx, err
	}
	return ctx, exec.Command("open", temp).Run()
}

type Node struct {
	Id       string `json:"id"`
	FileName string `json:"fileName"`
	DirName  string `json:"dirName"`
	Loc      int    `json:"loc"`
	Size     int    `json:"size"`
}

type Link struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

func marshalGraph(g *language.Graph, parser language.NodeParser) ([]byte, error) {
	out := Graph{}

	allNodes := g.AllNodes()
	maxLoc := max(utils.Max(allNodes, func(n *language.Node) int {
		return n.Data.Loc
	}), 1)

	for _, node := range allNodes {
		filename := parser.Display(node)
		out.Nodes = append(out.Nodes, Node{
			Id:       node.Id,
			FileName: path.Base(filename),
			DirName:  path.Dir(filename) + "/",
			Loc:      node.Data.Loc,
			Size:     10 * node.Data.Loc / maxLoc,
		})

		for _, to := range g.FromId(node.Id) {
			out.Links = append(out.Links, Link{
				From: node.Id,
				To:   to.Id,
			})
		}
	}

	return json.Marshal(out)
}
