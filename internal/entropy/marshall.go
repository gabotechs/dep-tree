package entropy

import (
	"encoding/json"
	"path"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/utils"
)

const (
	maxNodeSize         = 10
	folderNodeRepulsion = 1
	folderLinkStrength  = 2
	cyclicDepColor      = "indianred"
)

type Node struct {
	Id       string `json:"id"`
	FileName string `json:"fileName"`
	DirName  string `json:"dirName"`
	Loc      int    `json:"loc"`
	Size     int    `json:"size"`
	Color    []int  `json:"color,omitempty"`
	Visible  bool   `json:"visible"`
	// Repulsion factor to multiply the default repulsion that happens between nodes.
	Repulsion float64 `json:"repulsion,omitempty"`
}

type Link struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Color   string `json:"color,omitempty"`
	Visible bool   `json:"visible"`
	// TODO: I don't know what Distance means really, or how does it differ from Strength
	// Distance factor to multiply the default distance of an edge between nodes.
	Distance float64 `json:"distance,omitempty"`
	// Strength factor to multiply the default strength of a link.
	Strength int `json:"strength,omitempty"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

func marshal(dt *dep_tree.DepTree[language.FileInfo], parser language.NodeParser) ([]byte, error) {
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
			Visible:  true,
		})

		for _, to := range dt.Graph.FromId(node.Id) {
			out.Links = append(out.Links, Link{
				From:    node.Id,
				To:      to.Id,
				Visible: true,
			})
		}

		for _, parentFolder := range splitFullPaths(dirName) {
			out.Links = append(out.Links, Link{
				From:     node.Id,
				To:       parentFolder,
				Visible:  false,
				Strength: folderLinkStrength,
			})
			if _, ok := addedFolders[parentFolder]; ok {
				continue
			} else {
				addedFolders[parentFolder] = true
				out.Nodes = append(out.Nodes, Node{
					Id:        parentFolder,
					Visible:   false,
					Repulsion: folderNodeRepulsion,
				})
			}
		}
	}

	for el := dt.Cycles.Front(); el != nil; el = el.Next() {
		out.Links = append(out.Links, Link{
			From:    el.Key[0],
			To:      el.Key[1],
			Color:   cyclicDepColor,
			Visible: true,
		})
	}

	return json.Marshal(out)
}
