package entropy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/utils"
)

const (
	maxNodeSize = 10
)

type Node struct {
	Id           int64  `json:"id"`
	IsEntrypoint bool   `json:"isEntrypoint"`
	FileName     string `json:"fileName"`
	Group        string `json:"group,omitempty"`
	DirName      string `json:"dirName"`
	Loc          int    `json:"loc"`
	Size         int    `json:"size"`
	Color        []int  `json:"color,omitempty"`
	IsDir        bool   `json:"isDir"`
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

func makeGraph(files []string, parser graph.NodeParser[*language.FileInfo], loadCallbacks graph.LoadCallbacks[*language.FileInfo]) (Graph, error) {
	g := graph.NewGraph[*language.FileInfo]()
	err := g.Load(files, parser, loadCallbacks)
	if err != nil {
		return Graph{}, err
	}
	var singleEntrypointAbsPath string
	var entrypoints []*graph.Node[*language.FileInfo]
	if len(files) == 1 {
		entrypoint := g.Get(files[0])
		if entrypoint == nil {
			return Graph{}, fmt.Errorf("could not find entrypoint %s", files[0])
		}
		entrypoints = append(entrypoints, entrypoint)
		singleEntrypointAbsPath = entrypoint.Data.AbsPath
	} else {
		entrypoints = g.GetNodesWithoutParents()
	}

	cycles := g.RemoveCycles(entrypoints)
	out := Graph{
		Nodes: make([]Node, 0),
		Links: make([]Link, 0),
	}

	allNodes := g.AllNodes()
	maxLoc := max(utils.Max(allNodes, func(n *graph.Node[*language.FileInfo]) int {
		return n.Data.Loc
	}), 1)

	dirTree := NewDirTree()

	for _, node := range allNodes {
		dirTree.AddDirsFromFileInfo(node.Data)
	}

	addedFolders := map[string]bool{}

	for _, node := range allNodes {
		dirName := filepath.Dir(node.Data.RelPath)
		n := Node{
			Id:           node.ID(),
			IsEntrypoint: node.Data.AbsPath == singleEntrypointAbsPath,
			FileName:     filepath.Base(node.Data.RelPath),
			Group:        node.Data.Package,
			DirName:      dirName + string(os.PathSeparator),
			Loc:          node.Data.Loc,
			Size:         maxNodeSize * node.Data.Loc / maxLoc,
			Color:        toInt(dirTree.ColorForFileInfo(node.Data, RGB)),
		}
		out.Nodes = append(out.Nodes, n)

		for _, to := range g.FromId(node.Id) {
			out.Links = append(out.Links, Link{
				From: node.ID(),
				To:   to.ID(),
			})
		}

		for _, parentFolder := range dirTree.GroupingsForFileInfo(node.Data) {
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

	for _, cycle := range cycles {
		out.Links = append(out.Links, Link{
			From:     graph.MakeNode(cycle.Cause[0], 0).ID(),
			To:       graph.MakeNode(cycle.Cause[1], 0).ID(),
			IsCyclic: true,
		})
	}

	return out, nil
}
