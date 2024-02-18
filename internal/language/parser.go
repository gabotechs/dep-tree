package language

import (
	"github.com/elliotchance/orderedmap/v2"
	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/utils"
)

type Parser struct {
	Lang               Language
	UnwrapProxyExports bool
	Exclude            []string
	// cache
	FileCache    map[string]*FileInfo
	ImportsCache map[string]*ImportsResult
	ExportsCache map[string]*ExportEntries
}

func NewParser(lang Language) *Parser {
	return &Parser{
		Lang:               lang,
		UnwrapProxyExports: false,
		Exclude:            nil,
		FileCache:          make(map[string]*FileInfo),
		ImportsCache:       make(map[string]*ImportsResult),
		ExportsCache:       make(map[string]*ExportEntries),
	}
}

var _ graph.NodeParser[*FileInfo] = &Parser{}

func (p *Parser) shouldExclude(path string) bool {
	for _, exclusion := range p.Exclude {
		if ok, _ := utils.GlobstarMatch(exclusion, path); ok {
			return true
		}
	}
	return false
}

func (p *Parser) Node(id string) (*graph.Node[*FileInfo], error) {
	file, err := p.parseFile(id)
	if err != nil {
		return nil, err
	}
	return graph.MakeNode(id, file), nil
}

func (p *Parser) Deps(n *graph.Node[*FileInfo]) ([]*graph.Node[*FileInfo], error) {
	imports, err := p.gatherImportsFromFile(n.Id)
	if err != nil {
		return nil, err
	}
	n.AddErrors(imports.Errors...)

	// Some exports might be re-exporting symbols from other files, we consider
	// those as if they were normal imports.
	//
	// NOTE: if exports are parsed as imports, they might say that a name is being
	//  imported from a path when it's actually not available.
	//  ex:
	//   index.ts -> import { foo } from 'foo.ts'
	//   foo.ts   -> import { bar as foo } from 'bar.ts'
	//   bar.ts   -> export { bar }
	//  If unwrappedExports is true, this will say that `foo` is exported from `bar.ts`, which
	//  technically is true, but it's not true to say that `foo` is imported from `bar.ts`.
	//  It's more accurate to say that `bar` is imported from `bar.ts`, even if the alias is `foo`.
	//  Instead, we never unwrap export to avoid this.
	exports, err := p.parseExports(n.Id, false, nil)
	if err != nil {
		return nil, err
	}
	n.AddErrors(exports.Errors...)
	for el := exports.Symbols.Front(); el != nil; el = el.Next() {
		if el.Value != n.Id {
			imports.Imports = append(imports.Imports, ImportEntry{
				Symbols: []string{el.Key},
				AbsPath: el.Value,
			})
		}
	}

	resolvedImports := orderedmap.NewOrderedMap[string, bool]()

	// Imported names might not necessarily be declared in the path that is being imported, they might be declared in
	// a different file, we want that file. Ex: foo.ts -> utils/index.ts -> utils/sum.ts. If unwrapProxyExports is
	// set to true, we must trace those exports back.
	for _, importEntry := range imports.Imports {
		if !p.UnwrapProxyExports {
			resolvedImports.Set(importEntry.AbsPath, true)
			continue
		}

		// NOTE: at this point p.unwrapProxyExports is always true.
		exports, err = p.parseExports(importEntry.AbsPath, p.UnwrapProxyExports, nil)
		if err != nil {
			return nil, err
		}
		n.AddErrors(exports.Errors...)
		if importEntry.All {
			// If all imported, then dump every path in the resolved imports.
			for el := exports.Symbols.Front(); el != nil; el = el.Next() {
				resolvedImports.Set(el.Value, true)
			}
		} else if len(importEntry.Symbols) == 0 {
			resolvedImports.Set(importEntry.AbsPath, true)
		} else {
			for _, name := range importEntry.Symbols {
				if exportPath, ok := exports.Symbols.Get(name); ok {
					resolvedImports.Set(exportPath, true)
				} else {
					// TODO: this is not retro-compatible, do it in a different PR.
					// n.AddErrors(fmt.Errorf("name %s is imported by %s but not exported by %s", name, n.Id, importEntry.Id)).
				}
			}
		}
	}

	deps := make([]*graph.Node[*FileInfo], 0)
	for _, imported := range resolvedImports.Keys() {
		if p.shouldExclude(imported) {
			continue
		}
		node, err := p.Node(imported)
		if err != nil {
			n.AddErrors(err)
			continue
		}
		deps = append(deps, node)
	}
	return deps, nil
}
