package language

import (
	"context"
	"path"
	"path/filepath"

	"github.com/elliotchance/orderedmap/v2"

	"dep-tree/internal/dep_tree"
	"dep-tree/internal/graph"
)

type Language[T any, F any] interface {
	ParseFile(id string) (*F, error)
	MakeNode(id string) (*graph.Node[T], error)
	ParseImports(file *F) (*ImportsResult, error)
	ParseExports(file *F) (*ExportsResult, error)
}

type Parser[T any, F any] struct {
	entrypoint *graph.Node[T]
	lang       Language[T, F]
}

var _ dep_tree.NodeParser[any] = &Parser[any, any]{}

func makeParser[T any, F any](entrypoint string, languageBuilder func(string) (Language[T, F], error)) (*Parser[T, F], error) {
	lang, err := languageBuilder(entrypoint)
	if err != nil {
		return nil, err
	}
	entrypointNode, err := lang.MakeNode(entrypoint)
	return &Parser[T, F]{
		entrypoint: entrypointNode,
		lang:       lang,
	}, err
}

type Builder[T any, F any] func(string) (Language[T, F], error)

func ParserBuilder[T any, F any](languageBuilder Builder[T, F]) func(string) (dep_tree.NodeParser[T], error) {
	return func(entrypoint string) (dep_tree.NodeParser[T], error) {
		return makeParser[T, F](entrypoint, languageBuilder)
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
		ctx, exports, err = p.CachedUnwrappedParseExports(ctx, importEntry.Id)
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
