package js

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"dep-tree/internal/graph"
	"dep-tree/internal/node"
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

type Parser struct {
	RootDir string
}

var _ graph.NodeParser[Data] = &Parser{}

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

func (p *Parser) Parse(id string) (*node.Node[Data], error) {
	absPath, err := filepath.Abs(id)
	if err != nil {
		return nil, err
	}
	stat, _ := os.Stat(id)
	if stat != nil && stat.IsDir() {
		newAbsPath := retrieveWithExt(path.Join(absPath, "index"))
		if newAbsPath == "" {
			return nil, fmt.Errorf("tried to import from dir %s, but there is no index file", absPath)
		}
		absPath = newAbsPath
	} else {
		newAbsPath := retrieveWithExt(absPath)
		if newAbsPath == "" {
			return nil, fmt.Errorf("no matching JS externsion for file %s", absPath)
		}
		absPath = newAbsPath
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

func (p *Parser) Deps(n *node.Node[Data]) []string {
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

func (p *Parser) Display(n *node.Node[Data]) string {
	return path.Join(path.Base(path.Dir(n.Id)), path.Base(n.Id))
}
