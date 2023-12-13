package language

import (
	"context"
	"path"
	"path/filepath"

	"github.com/elliotchance/orderedmap/v2"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/graph"
)

type Language[T any, F any] interface {
	// ParseFile receives an absolute file path and returns F, where F is the specific file implementation
	//  defined by the language. This file object F will be used as input for parsing imports and exports.
	ParseFile(path string) (*F, error)
	// MakeNode receives an absolute file path and returns a graph.Node implementation.
	MakeNode(path string) (*graph.Node[T], error)
	// ParseImports receives the file F parsed by the ParseFile method and gathers the imports that the file
	//  F contains.
	ParseImports(file *F) (*ImportsResult, error)
	// ParseExports receives the file F parsed by the ParseFile method and gathers the exports that the file
	//  F contains.
	ParseExports(file *F) (*ExportsResult, error)
}

type Parser[T any, F any] struct {
	entrypoint *graph.Node[T]
	lang       Language[T, F]
}

var _ dep_tree.NodeParser[any] = &Parser[any, any]{}

func makeParser[T any, F any](ctx context.Context, entrypoint string, languageBuilder Builder[T, F]) (context.Context, *Parser[T, F], error) {
	ctx, lang, err := languageBuilder(ctx, entrypoint)
	if err != nil {
		return ctx, nil, err
	}
	entrypointNode, err := lang.MakeNode(entrypoint)
	return ctx, &Parser[T, F]{
		entrypoint: entrypointNode,
		lang:       lang,
	}, err
}

type Builder[T any, F any] func(context.Context, string) (context.Context, Language[T, F], error)

func ParserBuilder[T any, F any](languageBuilder Builder[T, F]) dep_tree.NodeParserBuilder[T] {
	return func(ctx context.Context, entrypoint string) (context.Context, dep_tree.NodeParser[T], error) {
		return makeParser[T, F](ctx, entrypoint, languageBuilder)
	}
}

func (p *Parser[T, F]) Entrypoint() (*graph.Node[T], error) {
	return p.entrypoint, nil
}

func (p *Parser[T, F]) Deps(ctx context.Context, n *graph.Node[T]) (context.Context, []*graph.Node[T], error) {
	ctx, imports, err := p.CachedParseImports(ctx, n.Id)
	if err != nil {
		return ctx, nil, err
	}
	n.AddErrors(imports.Errors...)

	resolvedImports := orderedmap.NewOrderedMap[string, bool]()

	// Take exports into account if top level root node is exporting stuff.
	if n.Id == p.entrypoint.Id {
		var exports *UnwrappedExportsResult
		ctx, exports, err = p.CachedUnwrappedParseExports(ctx, n.Id)
		if err != nil {
			return nil, nil, err
		}
		n.AddErrors(exports.Errors...)
		for _, k := range exports.Exports.Keys() {
			exportFrom, _ := exports.Exports.Get(k)
			resolvedImports.Set(exportFrom, true)
		}
	}

	// Imported names might not necessarily be declared in the path that is being imported, they might be declared in
	// a different file, we want that file. Ex: foo.ts -> utils/index.ts -> utils/sum.ts.
	for _, importEntry := range imports.Imports {
		var exports *UnwrappedExportsResult
		ctx, exports, err = p.CachedUnwrappedParseExports(ctx, importEntry.Path)
		if err != nil {
			return ctx, nil, err
		}
		n.AddErrors(exports.Errors...)
		if importEntry.All {
			// If all imported, then dump every path in the resolved imports.
			for _, k := range exports.Exports.Keys() {
				fromPath, _ := exports.Exports.Get(k)
				if _, ok := resolvedImports.Get(fromPath); ok {
					continue
				}
				resolvedImports.Set(fromPath, true)
			}
		} else {
			for _, name := range importEntry.Names {
				if resolvedImport, ok := exports.Exports.Get(name); ok {
					if _, ok := resolvedImports.Get(resolvedImport); ok {
						continue
					}
					resolvedImports.Set(resolvedImport, true)
				} else {
					// TODO: this is not retro-compatible, do it in a different PR.
					// n.AddErrors(fmt.Errorf("name %s is imported by %s but not exported by %s", name, n.Id, importEntry.Id)).
				}
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

func (p *Parser[T, F]) Display(n *graph.Node[T]) string {
	base := path.Dir(p.entrypoint.Id)
	rel, err := filepath.Rel(base, n.Id)
	if err != nil {
		return n.Id
	} else {
		return rel
	}
}
