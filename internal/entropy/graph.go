package entropy

import (
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/utils"
)

const (
	maxNodeSize = 10
)

type Node struct {
	Id       int64  `json:"id"`
	FileName string `json:"fileName"`
	Group    string `json:"group,omitempty"`
	DirName  string `json:"dirName"`
	Loc      int    `json:"loc"`
	Size     int    `json:"size"`
	Color    []int  `json:"color,omitempty"`
	IsDir    bool   `json:"isDir"`
}

type Link struct {
	From     int64 `json:"from"`
	To       int64 `json:"to"`
	IsDir    bool  `json:"isDir"`
	IsCyclic bool  `json:"isCyclic"`
}

type Graph struct {
	Nodes     []Node `json:"nodes"`
	Links     []Link `json:"links"`
	EnableGui bool   `json:"enableGui"`
}

func toInt(arr []float64) []int {
	result := make([]int, len(arr))
	for i, v := range arr {
		result[i] = int(v)
	}
	return result
}

func makeGraph(dt *dep_tree.DepTree[language.FileInfo], parser language.NodeParser) Graph {
	out := Graph{
		Nodes: make([]Node, 0),
		Links: make([]Link, 0),
	}

	allNodes := dt.Graph.AllNodes()
	maxLoc := max(utils.Max(allNodes, func(n *language.Node) int {
		return n.Data.Loc
	}), 1)

	addedFolders := map[string]bool{}

	dirTree := NewDirTree()

	for _, node := range allNodes {
		dirTree.AddDirsFromDisplay(parser.Display(node))
	}

	for _, node := range allNodes {
		display := parser.Display(node)
		dirName := filepath.Dir(display.Name)
		out.Nodes = append(out.Nodes, Node{
			Id:       node.ID(),
			FileName: filepath.Base(display.Name),
			Group:    display.Group,
			DirName:  dirName + "/",
			Loc:      node.Data.Loc,
			Size:     maxNodeSize * node.Data.Loc / maxLoc,
			Color:    toInt(dirTree.ColorForDisplay(display, RGB)),
		})

		for _, to := range dt.Graph.FromId(node.Id) {
			out.Links = append(out.Links, Link{
				From: node.ID(),
				To:   to.ID(),
			})
		}

		for _, parentFolder := range dirTree.GroupingsForDisplay(display) {
			folderNode := graph.MakeNode(parentFolder, 0)
			out.Links = append(out.Links, Link{
				From:  node.ID(),
				To:    folderNode.ID(),
				IsDir: true,
			})
			if _, ok := addedFolders[parentFolder]; ok {
				continue
			} else {
				addedFolders[parentFolder] = true
				out.Nodes = append(out.Nodes, Node{
					Id:    folderNode.ID(),
					IsDir: true,
				})
			}
		}
	}

	for el := dt.Cycles.Front(); el != nil; el = el.Next() {
		out.Links = append(out.Links, Link{
			From:     graph.MakeNode(el.Key[0], 0).ID(),
			To:       graph.MakeNode(el.Key[1], 0).ID(),
			IsCyclic: true,
		})
	}

	return out
}
