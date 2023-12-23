package entropy

import (
	"path"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/utils"
)

const (
	maxNodeSize = 10
)

type Node struct {
	Id       string `json:"id"`
	FileName string `json:"fileName"`
	DirName  string `json:"dirName"`
	Loc      int    `json:"loc"`
	Size     int    `json:"size"`
	Color    []int  `json:"color,omitempty"`
	IsDir    bool   `json:"isDir"`
}

type Link struct {
	From     string `json:"from"`
	To       string `json:"to"`
	IsDir    bool   `json:"isDir"`
	IsCyclic bool   `json:"isCyclic"`
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
		dirTree.AddDirs(path.Dir(parser.Display(node)))
	}

	for _, node := range allNodes {
		filepath := parser.Display(node)
		dirName := path.Dir(filepath)
		out.Nodes = append(out.Nodes, Node{
			Id:       node.Id,
			FileName: path.Base(filepath),
			DirName:  dirName + "/",
			Loc:      node.Data.Loc,
			Size:     maxNodeSize * node.Data.Loc / maxLoc,
			Color:    dirTree.ColorFor(dirName),
		})

		for _, to := range dt.Graph.FromId(node.Id) {
			out.Links = append(out.Links, Link{
				From: node.Id,
				To:   to.Id,
			})
		}

		for _, parentFolder := range splitFullPaths(dirName) {
			out.Links = append(out.Links, Link{
				From:  node.Id,
				To:    parentFolder,
				IsDir: true,
			})
			if _, ok := addedFolders[parentFolder]; ok {
				continue
			} else {
				addedFolders[parentFolder] = true
				out.Nodes = append(out.Nodes, Node{
					Id:    parentFolder,
					IsDir: true,
				})
			}
		}
	}

	for el := dt.Cycles.Front(); el != nil; el = el.Next() {
		out.Links = append(out.Links, Link{
			From:     el.Key[0],
			To:       el.Key[1],
			IsCyclic: true,
		})
	}

	return out
}
