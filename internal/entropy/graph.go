package entropy

import (
	"path/filepath"
	"strings"

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
		dirTree.AddDirs(filepath.Dir(parser.Display(node)))
	}

	for _, node := range allNodes {
		filePath := parser.Display(node)
		dirName := filepath.Dir(filePath)
		out.Nodes = append(out.Nodes, Node{
			Id:       node.ID(),
			FileName: filepath.Base(filePath),
			DirName:  dirName + "/",
			Loc:      node.Data.Loc,
			Size:     maxNodeSize * node.Data.Loc / maxLoc,
			Color:    dirTree.ColorFor(dirName),
		})

		for _, to := range dt.Graph.FromId(node.Id) {
			out.Links = append(out.Links, Link{
				From: node.ID(),
				To:   to.ID(),
			})
		}

		for _, parentFolder := range splitFullPaths(dirName) {
			// NOTE: just ignore parent folders like ".." or "../..", otherwise they will contribute
			//  to grouping folders that might be unrelated. Empirically, visualizations look nicer if
			//  we ignore them.
			if strings.HasSuffix(parentFolder, "..") {
				continue
			}
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
