package js

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"dep-tree/internal/graph"
	"dep-tree/internal/graph/node"
	"dep-tree/internal/utils"
)

type Parser struct {
	ProjectRoot string
	TsConfig    TsConfig
}

var _ graph.NodeParser[Data] = &Parser{}

func MakeJsParser(entrypoint string) (*Parser, error) {
	entrypointAbsPath, err := filepath.Abs(entrypoint)
	if err != nil {
		return nil, err
	}
	searchPath := entrypointAbsPath
	var tsConfig TsConfig
	var projectRoot string
	for len(searchPath) > 1 {
		packageJsonPath := path.Join(searchPath, "package.json")
		if utils.FileExists(packageJsonPath) {
			tsConfigPath := path.Join(searchPath, "tsConfig.json")
			if utils.FileExists(tsConfigPath) {
				var err error
				tsConfig, err = ParseTsConfig(tsConfigPath)
				if err != nil {
					return nil, fmt.Errorf("found TypeScript config file in %s but there was an error reading it: %w", tsConfigPath, err)
				}
			}
			projectRoot = searchPath
			break
		} else {
			searchPath = path.Dir(searchPath)
		}
	}
	return &Parser{
		ProjectRoot: projectRoot,
		TsConfig:    tsConfig,
	}, nil
}

func (p *Parser) Entrypoint(entrypoint string) (*node.Node[Data], error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	} else if !utils.FileExists(entrypoint) {
		return nil, fmt.Errorf("file '%s' does not exist or is not visible form CWD %s", entrypoint, cwd)
	}
	entrypointAbsPath, err := filepath.Abs(entrypoint)
	if err != nil {
		return nil, err
	}

	return MakeJsNode(entrypointAbsPath)
}

func (p *Parser) Deps(n *node.Node[Data]) ([]*node.Node[Data], error) {
	fileInfo, err := p.ParseFileInfo(n.Data.content, n.Data.dirname)
	if err != nil {
		return nil, err
	}

	deps := make([]*node.Node[Data], 0)
	for _, imported := range fileInfo.imports {
		dep, err := MakeJsNode(imported.AbsPath)
		if err != nil {
			return nil, err
		}
		deps = append(deps, dep)
	}
	return deps, nil
}

func (p *Parser) Display(n *node.Node[Data], root *node.Node[Data]) string {
	base := root.Data.dirname
	rel, err := filepath.Rel(base, n.Id)
	if err != nil {
		return n.Id
	} else {
		return rel
	}
}
