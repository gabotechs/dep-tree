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
	"github.com/gabotechs/dep-tree/internal/graph"
)

//go:embed index.html
var index string

const ToReplace = "const GRAPH = {}"
const ReplacePrefix = "const GRAPH = "

func Render[T any](ctx context.Context, parser dep_tree.NodeParser[T]) (context.Context, error) {
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
	// Node identifier.
	Id string `json:"id"`
	// Node object accessor function or attribute for name (shown in label).
	// Supports plain text or HTML content. Note that this method uses innerHTML internally,
	// so make sure to pre-sanitize any user-input content to prevent XSS vulnerabilities.
	name string `json:"name"`
	// Node object accessor function, attribute or a numeric constant for the node
	// numeric value (affects sphere volume).
	val int `json:"val"`
	// Node object accessor function or attribute for node color (affects sphere color)a .
	Color string `json:"color,omitempty"`
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Color  string `json:"color,omitempty"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

func marshalGraph[T any](g *graph.Graph[T], parser dep_tree.NodeParser[T]) ([]byte, error) {
	out := Graph{}
	for _, node := range g.AllNodes() {
		out.Nodes = append(out.Nodes, Node{
			Id:   node.Id,
			name: parser.Display(node),
			val:  10,
		})

		for _, edge := range g.FromId(node.Id) {
			out.Links = append(out.Links, Link{
				Source: node.Id,
				Target: edge.Id,
			})
		}
	}

	return json.Marshal(out)
}
