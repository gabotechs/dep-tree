package js

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"dep-tree/internal/graph"
	"dep-tree/internal/graph/node"
)

var Extensions = []string{
	"js", "ts", "tsx", "jsx",
}

var importRegex = regexp.MustCompile(
	"(import|export)\\s+?((([\\w*\\s{},]*)\\s+from\\s+?)|)((\".*?\")|('.*?'))\\s*?(?:;|$|)",
)

var importPathRegex = regexp.MustCompile(
	"(\".*?\")|('.*?')",
)

type Data struct {
	dirname string
	content []byte
}

type parser struct{}

var Parser graph.NodeParser[Data] = &parser{}

func retrieveWithExt(absPath string) string {
	for _, ext := range Extensions {
		if strings.HasSuffix(absPath, "."+ext) {
			return absPath
		}
	}
	for _, ext := range Extensions {
		withExtPath := absPath + "." + ext
		_, err := os.Stat(withExtPath)
		if err == nil {
			return withExtPath
		}
	}
	return ""
}

func normalizeId(id string) (string, error) {
	absPath, err := filepath.Abs(id)
	if err != nil {
		return "", err
	}
	stat, _ := os.Stat(id)
	if stat != nil && stat.IsDir() {
		newAbsPath := retrieveWithExt(path.Join(absPath, "index"))
		if newAbsPath == "" {
			return "", fmt.Errorf("tried to import from dir %s, but there is no index file", absPath)
		}
		absPath = newAbsPath
	} else {
		newAbsPath := retrieveWithExt(absPath)
		if newAbsPath == "" {
			return "", fmt.Errorf("no matching JS extension for file %s", absPath)
		}
		absPath = newAbsPath
	}
	return absPath, nil
}

func (p *parser) Parse(id string) (*node.Node[Data], error) {
	absPath, err := normalizeId(id)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	dirname := path.Dir(absPath)

	return node.MakeNode(absPath, dirname, Data{
		dirname: dirname,
		content: content,
	}), nil
}

func (p *parser) Deps(n *node.Node[Data]) []string {
	matched := importRegex.FindAll(n.Data.content, -1)
	deps := make([]string, 0)
	for _, importMatch := range matched {
		importPathMatched := importPathRegex.Find(importMatch)
		match := strings.Trim(string(importPathMatched), "'\" \n")
		if match[:1] != "." {
			continue
		}
		match = path.Join(n.Data.dirname, match)
		deps = append(deps, match)
	}
	return deps
}

func (p *parser) Display(n *node.Node[Data]) string {
	return path.Join(path.Base(path.Dir(n.Id)), path.Base(n.Id))
}
