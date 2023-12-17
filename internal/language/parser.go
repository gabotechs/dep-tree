package language

import (
	"context"
	"path"
	"path/filepath"

	"github.com/elliotchance/orderedmap/v2"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/graph"
)

type Node = graph.Node[FileInfo]
type Graph = graph.Graph[FileInfo]
type NodeParser = dep_tree.NodeParser[FileInfo]
type NodeParserBuilder = dep_tree.NodeParserBuilder[FileInfo]

type FileInfo struct {
	Loc  int
	Size int
}

type CodeFile interface {
	Loc() int
	Size() int
}

type Language[F CodeFile] interface {
	// ParseFile receives an absolute file path and returns F, where F is the specific file implementation
	//  defined by the language. This file object F will be used as input for parsing imports and exports.
	ParseFile(path string) (*F, error)
	// ParseImports receives the file F parsed by the ParseFile method and gathers the imports that the file
	//  F contains.
	ParseImports(file *F) (*ImportsResult, error)
	// ParseExports receives the file F parsed by the ParseFile method and gathers the exports that the file
	//  F contains.
	ParseExports(file *F) (*ExportsResult, error)
}

type Parser[F CodeFile] struct {
	entrypoint *Node
	lang       Language[F]
}

var _ NodeParser = &Parser[CodeFile]{}

func makeParser[F CodeFile, C any](ctx context.Context, entrypoint string, languageBuilder Builder[F, C], cfg C) (context.Context, *Parser[F], error) {
	ctx, lang, err := languageBuilder(ctx, entrypoint, cfg)
	if err != nil {
		return ctx, nil, err
	}
	absEntrypoint, err := filepath.Abs(entrypoint)
	if err != nil {
		return ctx, nil, err
	}

	entrypointNode := graph.MakeNode(absEntrypoint, FileInfo{})
	return ctx, &Parser[F]{
		entrypoint: entrypointNode,
		lang:       lang,
	}, err
}

type Builder[F CodeFile, C any] func(context.Context, string, C) (context.Context, Language[F], error)

func ParserBuilder[F CodeFile, C any](languageBuilder Builder[F, C], cfg C) NodeParserBuilder {
	return func(ctx context.Context, entrypoint string) (context.Context, NodeParser, error) {
		return makeParser[F](ctx, entrypoint, languageBuilder, cfg)
	}
}

func (p *Parser[F]) Entrypoint() (*Node, error) {
	return p.entrypoint, nil
}

func (p *Parser[F]) Deps(ctx context.Context, n *Node) (context.Context, []*Node, error) {
	ctx, file, err := p.CachedParseFile(ctx, n.Id)
	if err != nil {
		return ctx, nil, err
	}
	n.Data.Size = (*file).Size()
	n.Data.Loc = (*file).Loc()
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

	deps := make([]*Node, resolvedImports.Len())
	for i, imported := range resolvedImports.Keys() {
		deps[i] = graph.MakeNode(imported, FileInfo{})
	}
	return ctx, deps, nil
}

func (p *Parser[F]) Display(n *Node) string {
	base := path.Dir(p.entrypoint.Id)
	rel, err := filepath.Rel(base, n.Id)
	if err != nil {
		return n.Id
	} else {
		return rel
	}
}
