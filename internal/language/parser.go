package language

import (
	"context"
	"path"
	"path/filepath"

	"github.com/elliotchance/orderedmap/v2"

	"dep-tree/internal/dep_tree"
	"dep-tree/internal/graph"
)

type Language[T any] interface {
	MakeNode(id string) (*graph.Node[T], error)
	ParseImports(ctx context.Context, id string) (context.Context, *ImportsResult, error)
	ParseExports(ctx context.Context, id string) (context.Context, *ExportsResult, error)
}

type Parser[T any] struct {
	entrypoint *graph.Node[T]
	lang       Language[T]
}

var _ dep_tree.NodeParser[any] = &Parser[any]{}

func makeParser[T any](entrypoint string, languageBuilder func(string) (Language[T], error)) (*Parser[T], error) {
	lang, err := languageBuilder(entrypoint)
	if err != nil {
		return nil, err
	}
	entrypointNode, err := lang.MakeNode(entrypoint)
	return &Parser[T]{
		entrypoint: entrypointNode,
		lang:       lang,
	}, err
}

func ParserBuilder[T any](languageBuilder func(string) (Language[T], error)) func(string) (dep_tree.NodeParser[T], error) {
	return func(entrypoint string) (dep_tree.NodeParser[T], error) {
		return makeParser[T](entrypoint, languageBuilder)
	}
}

func (p *Parser[T]) Entrypoint() (*graph.Node[T], error) {
	return p.entrypoint, nil
}

func (p *Parser[T]) Deps(ctx context.Context, n *graph.Node[T]) (context.Context, []*graph.Node[T], error) {
	ctx, imports, err := p.lang.ParseImports(ctx, n.Id)
	if err != nil {
		return ctx, nil, err
	}
	n.AddErrors(imports.Errors...)

	resolvedImports := orderedmap.NewOrderedMap[string, bool]()

	// Take exports into account if top level root node is exporting stuff.
	if n.Id == p.entrypoint.Id {
		var exports *ExportsResult
		ctx, exports, err = p.lang.ParseExports(ctx, n.Id)
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
		ctx, exports, err = p.lang.ParseExports(ctx, importedPath)
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

	deps := make([]*graph.Node[T], 0)
	for _, imported := range resolvedImports.Keys() {
		dep, err := p.lang.MakeNode(imported)
		if err != nil {
			return ctx, nil, err
		}
		deps = append(deps, dep)
	}
	return ctx, deps, nil
}

func (p *Parser[T]) Display(n *graph.Node[T]) string {
	base := path.Dir(p.entrypoint.Id)
	rel, err := filepath.Rel(base, n.Id)
	if err != nil {
		return n.Id
	} else {
		return rel
	}
}
