package js

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/elliotchance/orderedmap/v2"

	"dep-tree/internal/graph"
	"dep-tree/internal/utils"
)

type Parser struct {
	entrypoint  string
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
			tsConfigPath := path.Join(searchPath, "tsconfig.json")
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
		entrypoint:  entrypointAbsPath,
		ProjectRoot: projectRoot,
		TsConfig:    tsConfig,
	}, nil
}

func (p *Parser) Entrypoint(entrypoint string) (*graph.Node[Data], error) {
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

func (p *Parser) Deps(ctx context.Context, n *graph.Node[Data]) (context.Context, []*graph.Node[Data], error) {
	ctx, imports, err := p.parseImports(ctx, n.Data.filePath)
	if err != nil {
		return ctx, nil, err
	}

	// Imported names might not necessarily be declared in the path that is being imported, they might be declared in
	// a different file, we want that file. Ex: foo.ts -> utils/index.ts -> utils/sum.ts.
	resolvedImports := orderedmap.NewOrderedMap[string, bool]()
	for _, importedPath := range imports.Keys() {
		importedNames, _ := imports.Get(importedPath)
		var exports map[string]string
		ctx, exports, err = p.parseExports(ctx, importedPath)
		if err != nil {
			return ctx, nil, err
		}
		for _, name := range importedNames {
			// If all imported, then dump every path in the resolved imports.
			if name == "*" {
				for _, fromPath := range exports {
					if _, ok := resolvedImports.Get(fromPath); ok {
						continue
					}
					resolvedImports.Set(fromPath, true)
				}
				break
			}

			if resolvedImport, ok := exports[name]; ok {
				if _, ok := resolvedImports.Get(resolvedImport); ok {
					continue
				}
				resolvedImports.Set(resolvedImport, true)
			}
		}
	}

	deps := make([]*graph.Node[Data], 0)
	for _, imported := range resolvedImports.Keys() {
		dep, err := MakeJsNode(imported)
		if err != nil {
			return ctx, nil, err
		}
		deps = append(deps, dep)
	}
	return ctx, deps, nil
}

func (p *Parser) Display(n *graph.Node[Data]) string {
	base := path.Dir(p.entrypoint)
	rel, err := filepath.Rel(base, n.Id)
	if err != nil {
		return n.Id
	} else {
		return rel
	}
}
