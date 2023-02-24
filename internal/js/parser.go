package js

import (
	"context"
	"fmt"
	"path"
	"path/filepath"

	"github.com/elliotchance/orderedmap/v2"

	"dep-tree/internal/dep_tree"
	"dep-tree/internal/graph"
	"dep-tree/internal/utils"
)

type Parser struct {
	rootId      string
	ProjectRoot string
	TsConfig    TsConfig
}

var _ dep_tree.NodeParser[Data] = &Parser{}

func findPackageJson(searchPath string) (TsConfig, string, error) {
	if len(searchPath) < 2 {
		return TsConfig{}, "", nil
	}
	packageJsonPath := path.Join(searchPath, "package.json")
	if utils.FileExists(packageJsonPath) {
		tsConfigPath := path.Join(searchPath, "tsconfig.json")
		var tsConfig TsConfig
		var err error
		if utils.FileExists(tsConfigPath) {
			tsConfig, err = ParseTsConfig(tsConfigPath)
			if err != nil {
				err = fmt.Errorf("found TypeScript config file in %s but there was an error reading it: %w", tsConfigPath, err)
			}
		}
		return tsConfig, searchPath, err
	} else {
		return findPackageJson(path.Dir(searchPath))
	}
}

func MakeJsParser(entrypoint string) (dep_tree.NodeParser[Data], error) {
	entrypointAbsPath, err := filepath.Abs(entrypoint)
	if err != nil {
		return nil, err
	}
	tsConfig, packageJsonPath, err := findPackageJson(entrypointAbsPath)
	if err != nil {
		return nil, err
	}
	projectRoot := path.Dir(entrypointAbsPath)
	if packageJsonPath != "" {
		projectRoot = packageJsonPath
	}
	return &Parser{
		rootId:      entrypointAbsPath,
		ProjectRoot: projectRoot,
		TsConfig:    tsConfig,
	}, nil
}

func (p *Parser) Entrypoint() (*graph.Node[Data], error) {
	return MakeJsNode(p.rootId)
}

func (p *Parser) Deps(ctx context.Context, n *graph.Node[Data]) (context.Context, []*graph.Node[Data], error) {
	ctx, imports, err := p.parseImports(ctx, n.Id)
	if err != nil {
		return ctx, nil, err
	}
	n.AddErrors(imports.Errors...)

	resolvedImports := orderedmap.NewOrderedMap[string, bool]()

	// Take exports into account if top level root node is exporting stuff.
	if n.Id == p.rootId {
		var exports *ExportsResult
		ctx, exports, err = p.parseExports(ctx, p.rootId)
		if err != nil {
			return nil, nil, err
		}
		n.AddErrors(exports.Errors...)
		for _, exportFrom := range exports.Exports {
			resolvedImports.Set(exportFrom, true)
		}
	}

	// Imported names might not necessarily be declared in the path that is being imported, they might be declared in
	// a different file, we want that file. Ex: foo.ts -> utils/index.ts -> utils/sum.ts.
	for _, importedPath := range imports.Imports.Keys() {
		importedNames, _ := imports.Imports.Get(importedPath)
		var exports *ExportsResult
		ctx, exports, err = p.parseExports(ctx, importedPath)
		if err != nil {
			return ctx, nil, err
		}
		n.AddErrors(exports.Errors...)
		for _, name := range importedNames {
			// If all imported, then dump every path in the resolved imports.
			if name == "*" {
				for _, fromPath := range exports.Exports {
					if _, ok := resolvedImports.Get(fromPath); ok {
						continue
					}
					resolvedImports.Set(fromPath, true)
				}
				break
			}

			if resolvedImport, ok := exports.Exports[name]; ok {
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
	base := path.Dir(p.rootId)
	rel, err := filepath.Rel(base, n.Id)
	if err != nil {
		return n.Id
	} else {
		return rel
	}
}
