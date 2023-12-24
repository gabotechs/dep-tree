package language

import (
	"context"
	"path"
	"path/filepath"

	"github.com/elliotchance/orderedmap/v2"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/utils"
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
	ParseExports(file *F) (*ExportsEntries, error)
}

type Parser[F CodeFile] struct {
	entrypoint         *Node
	lang               Language[F]
	unwrapProxyExports bool
	exclude            []string
}

var _ NodeParser = &Parser[CodeFile]{}

func makeParser[F CodeFile, C any](
	ctx context.Context,
	entrypoint string,
	languageBuilder Builder[F, C],
	langCfg C,
	generalCfg Config,
) (context.Context, *Parser[F], error) {
	ctx, lang, err := languageBuilder(ctx, entrypoint, langCfg)
	if err != nil {
		return ctx, nil, err
	}
	absEntrypoint, err := filepath.Abs(entrypoint)
	if err != nil {
		return ctx, nil, err
	}

	entrypointNode := graph.MakeNode(absEntrypoint, FileInfo{})
	parser := &Parser[F]{
		entrypoint:         entrypointNode,
		lang:               lang,
		unwrapProxyExports: true,
	}
	if generalCfg != nil {
		parser.unwrapProxyExports = generalCfg.UnwrapProxyExports()
		parser.exclude = generalCfg.IgnoreFiles()
	}
	return ctx, parser, err
}

type Config interface {
	UnwrapProxyExports() bool
	IgnoreFiles() []string
}

type Builder[F CodeFile, C any] func(context.Context, string, C) (context.Context, Language[F], error)

func ParserBuilder[F CodeFile, C any](languageBuilder Builder[F, C], langCfg C, generalCfg Config) NodeParserBuilder {
	return func(ctx context.Context, entrypoint string) (context.Context, NodeParser, error) {
		return makeParser[F](ctx, entrypoint, languageBuilder, langCfg, generalCfg)
	}
}

func (p *Parser[F]) shouldExclude(path string) bool {
	for _, exclusion := range p.exclude {
		if ok, _ := utils.GlobstarMatch(exclusion, path); ok {
			return true
		}
	}
	return false
}

func (p *Parser[F]) Entrypoint() (*Node, error) {
	return p.entrypoint, nil
}

func (p *Parser[F]) updateNodeInfo(ctx context.Context, n *Node) (context.Context, error) {
	ctx, file, err := p.CachedParseFile(ctx, n.Id)
	if err != nil {
		return ctx, err
	}
	n.Data.Size = (*file).Size()
	n.Data.Loc = (*file).Loc()
	return ctx, err
}

func (p *Parser[F]) Deps(ctx context.Context, n *Node) (context.Context, []*Node, error) {
	ctx, _ = p.updateNodeInfo(ctx, n)
	ctx, imports, err := p.CachedParseImports(ctx, n.Id)
	if err != nil {
		return ctx, nil, err
	}
	n.AddErrors(imports.Errors...)

	// If exports are not going to be unwrapped, then we should take an export as if it
	// was importing names into the file. This might happen because of configured that way
	// or because it's the root file.
	if !p.unwrapProxyExports || n.Id == p.entrypoint.Id {
		var exports *ExportsResult
		// TODO: if exports are parsed as imports, they might say that that a name is being
		//  imported from a path when it's actually not available.
		//  ex:
		//   index.ts -> import { foo } from 'foo.ts'
		//   foo.ts   -> import { bar as foo } from 'bar.ts'
		//   bar.ts   -> export { bar }
		//  If unwrappedExports is true, this will say that `foo` is exported from `bar.ts`, which
		//  technically is true, but it's not true to say that `foo` is imported from `bar.ts`.
		//  It's more accurate to say that `bar` is imported from `bar.ts`, even if the alias is `foo`.
		//  Instead we never unwrap export to avoid this.
		ctx, exports, err = p.ParseExports(ctx, n.Id, false)
		if err != nil {
			return nil, nil, err
		}
		n.AddErrors(exports.Errors...)
		for el := exports.Exports.Front(); el != nil; el = el.Next() {
			imports.Imports = append(imports.Imports, ImportEntry{
				Names: []string{el.Key},
				Path:  el.Value,
			})
		}
	}

	resolvedImports := orderedmap.NewOrderedMap[string, bool]()

	// Imported names might not necessarily be declared in the path that is being imported, they might be declared in
	// a different file, we want that file. Ex: foo.ts -> utils/index.ts -> utils/sum.ts. If unwrapProxyExports is
	// set to true, we must trace those exports back.
	for _, importEntry := range imports.Imports {
		if !p.unwrapProxyExports {
			resolvedImports.Set(importEntry.Path, true)
			continue
		}

		var exports *ExportsResult
		ctx, exports, err = p.ParseExports(ctx, importEntry.Path, true)
		if err != nil {
			return ctx, nil, err
		}
		n.AddErrors(exports.Errors...)
		if importEntry.All {
			// If all imported, then dump every path in the resolved imports.
			for el := exports.Exports.Front(); el != nil; el = el.Next() {
				resolvedImports.Set(el.Value, true)
			}
			continue
		}
		for _, name := range importEntry.Names {
			if exportPath, ok := exports.Exports.Get(name); ok {
				resolvedImports.Set(exportPath, true)
			} else {
				// TODO: this is not retro-compatible, do it in a different PR.
				// n.AddErrors(fmt.Errorf("name %s is imported by %s but not exported by %s", name, n.Id, importEntry.Id)).
			}
		}
	}

	deps := make([]*Node, 0)
	for _, imported := range resolvedImports.Keys() {
		node := graph.MakeNode(imported, FileInfo{})
		if !p.shouldExclude(p.Display(node)) {
			deps = append(deps, node)
		}
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
