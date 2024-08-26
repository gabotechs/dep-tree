package entropy

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/utils"
)

const (
	maxNodeSize = 10
)

type Node struct {
	Id           int64    `json:"id"`
	IsEntrypoint bool     `json:"isEntrypoint"`
	IsDirectory  bool     `json:"isDirectory"`
	FileName     string   `json:"fileName"`
	PathBuf      []string `json:"pathBuf"`
	Group        string   `json:"group,omitempty"`
	DirName      string   `json:"dirName"`
	Loc          int      `json:"loc"`
	Size         int      `json:"size"`
}

type Link struct {
	From     int64 `json:"from"`
	To       int64 `json:"to"`
	IsCyclic bool  `json:"isCyclic"`
}

type Graph struct {
	Nodes     []Node `json:"nodes"`
	Links     []Link `json:"links"`
	EnableGui bool   `json:"enableGui"`
}

func makeGraph(files []string, parser graph.NodeParser[*language.FileInfo], loadCallbacks graph.LoadCallbacks[*language.FileInfo]) (Graph, error) {
	g := graph.NewGraph[*language.FileInfo]()
	err := g.Load(files, parser, loadCallbacks)
	if err != nil {
		return Graph{}, err
	}

	out := Graph{
		Nodes: make([]Node, 0),
		Links: make([]Link, 0),
	}

	dirNodes := make(map[string]int64)

	allNodes := g.AllNodes()
	maxLoc := max(utils.Max(allNodes, func(n *graph.Node[*language.FileInfo]) int {
		return n.Data.Loc
	}), 1)

	// First, create all directory nodes
	for _, node := range allNodes {
		dirPath := filepath.Dir(node.Data.AbsPath)
		for dirPath != "." && dirPath != "/" {
			if _, exists := dirNodes[dirPath]; !exists {
				dirNode := Node{
					Id:          int64(len(out.Nodes)),
					IsDirectory: true,
					FileName:    filepath.Base(dirPath),
					PathBuf:     strings.Split(dirPath, string(os.PathSeparator)),
					DirName:     filepath.Dir(dirPath) + string(os.PathSeparator),
				}
				out.Nodes = append(out.Nodes, dirNode)
				dirNodes[dirPath] = dirNode.Id
			}
			dirPath = filepath.Dir(dirPath)
		}
	}

	// Then, create file nodes and connect them to their parent directories
	for _, node := range allNodes {
		fileNode := Node{
			Id:           node.ID(),
			IsEntrypoint: node.Data.AbsPath == files[0],
			FileName:     filepath.Base(node.Data.RelPath),
			PathBuf:      strings.Split(node.Data.AbsPath, string(os.PathSeparator)),
			Group:        node.Data.Package,
			DirName:      filepath.Dir(node.Data.RelPath) + string(os.PathSeparator),
			Loc:          node.Data.Loc,
			Size:         maxNodeSize * node.Data.Loc / maxLoc,
		}
		out.Nodes = append(out.Nodes, fileNode)

		// Connect file to its parent directory
		parentDir := filepath.Dir(node.Data.AbsPath)
		if parentDirId, exists := dirNodes[parentDir]; exists {
			out.Links = append(out.Links, Link{
				From: parentDirId,
				To:   fileNode.Id,
			})
		}

		// Add original file dependencies
		for _, to := range g.FromId(node.Id) {
			out.Links = append(out.Links, Link{
				From: node.ID(),
				To:   to.ID(),
			})
		}
	}

	// Connect directories to their parent directories
	for dirPath, dirId := range dirNodes {
		parentDir := filepath.Dir(dirPath)
		if parentDirId, exists := dirNodes[parentDir]; exists && parentDir != dirPath {
			out.Links = append(out.Links, Link{
				From: parentDirId,
				To:   dirId,
			})
		}
	}

	return out, nil
}
