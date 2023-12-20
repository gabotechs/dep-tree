package entropy

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"path"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/utils"
)

//go:embed index.html
var index []byte

const ToReplace = "const GRAPH = {}"
const ReplacePrefix = "const GRAPH = "

func Render(ctx context.Context, parser language.NodeParser) (context.Context, error) {
	dt := dep_tree.NewDepTree(parser)
	ctx, err := dt.LoadGraph(ctx)
	if err != nil {
		return ctx, err
	}

	dt.LoadCycles()
	marshaled, err := marshalGraph(dt, parser)
	if err != nil {
		return ctx, err
	}
	rendered := bytes.ReplaceAll(index, []byte(ToReplace), append([]byte(ReplacePrefix), marshaled...))
	temp := path.Join(os.TempDir(), "index.html")
	err = os.WriteFile(temp, rendered, os.ModePerm)
	if err != nil {
		return ctx, err
	}
	return ctx, openInBrowser(temp)
}

type Node struct {
	Id       string `json:"id"`
	FileName string `json:"fileName"`
	DirName  string `json:"dirName"`
	Loc      int    `json:"loc"`
	Size     int    `json:"size"`
	Visible  bool   `json:"visible"`
}

type Link struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Color   string `json:"color,omitempty"`
	Visible bool   `json:"visible"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

func marshalGraph(dt *dep_tree.DepTree[language.FileInfo], parser language.NodeParser) ([]byte, error) {
	out := Graph{}

	allNodes := dt.Graph.AllNodes()
	maxLoc := max(utils.Max(allNodes, func(n *language.Node) int {
		return n.Data.Loc
	}), 1)

	addedFolders := map[string]bool{}

	dirTree := NewDirTree()

	for _, node := range allNodes {
		filepath := parser.Display(node)
		dirName := path.Dir(filepath)
		out.Nodes = append(out.Nodes, Node{
			Id:       node.Id,
			FileName: path.Base(filepath),
			DirName:  dirName + "/",
			Loc:      node.Data.Loc,
			Size:     10 * node.Data.Loc / maxLoc,
			Visible:  true,
		})

		for _, to := range dt.Graph.FromId(node.Id) {
			out.Links = append(out.Links, Link{
				From:    node.Id,
				To:      to.Id,
				Visible: true,
			})
		}

		for _, parentFolder := range dirTree.AddDirs(dirName) {
			out.Links = append(out.Links, Link{
				From:    node.Id,
				To:      parentFolder,
				Visible: false,
			})
			if _, ok := addedFolders[parentFolder]; ok {
				continue
			} else {
				addedFolders[parentFolder] = true
				out.Nodes = append(out.Nodes, Node{
					Id:      parentFolder,
					Visible: false,
				})
			}
		}
	}

	for el := dt.Cycles.Front(); el != nil; el = el.Next() {
		out.Links = append(out.Links, Link{
			From:  el.Key[0],
			To:    el.Key[1],
			Color: "red",
		})
	}

	return json.Marshal(out)
}
