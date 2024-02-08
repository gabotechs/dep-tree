package language

import (
	"github.com/elliotchance/orderedmap/v2"
	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/utils"
)

type FileInfo struct {
	Content any
	Path    string
	Loc     int
	Size    int
}

type Language interface {
	// ParseFile receives an absolute file path and returns F, where F is the specific file implementation
	//  defined by the language. This file object F will be used as input for parsing imports and exports.
	ParseFile(path string) (*FileInfo, error)
	// ParseImports receives the file F parsed by the ParseFile method and gathers the imports that the file
	//  F contains.
	ParseImports(file *FileInfo) (*ImportsResult, error)
	// ParseExports receives the file F parsed by the ParseFile method and gathers the exports that the file
	//  F contains.
	ParseExports(file *FileInfo) (*ExportsEntries, error)
	// Display takes an absolute path to a file and displays it nicely.
	Display(path string) graph.DisplayResult
}

type Parser struct {
	lang               Language
	unwrapProxyExports bool
	exclude            []string
	// cache
	fileCache    map[string]*FileInfo
	importsCache map[string]*ImportsResult
	exportsCache map[string]*ExportsResult
}

var _ graph.NodeParser[*FileInfo] = &Parser{}

type Config interface {
	UnwrapProxyExports() bool
	IgnoreFiles() []string
}

type Builder[C any] func(C) (Language, error)

func ParserBuilder[C any](
	languageBuilder Builder[C],
	langCfg C,
	generalCfg Config,
) graph.NodeParserBuilder[*FileInfo] {
	fileCache := map[string]*FileInfo{}
	importsCache := map[string]*ImportsResult{}
	exportsCache := map[string]*ExportsResult{}
	return func(files []string) (graph.NodeParser[*FileInfo], error) {
		lang, err := languageBuilder(langCfg)
		if err != nil {
			return nil, err
		}

		parser := &Parser{
			lang:               lang,
			unwrapProxyExports: true,
			fileCache:          fileCache,
			importsCache:       importsCache,
			exportsCache:       exportsCache,
		}
		if generalCfg != nil {
			parser.unwrapProxyExports = generalCfg.UnwrapProxyExports()
			parser.exclude = generalCfg.IgnoreFiles()
		}
		return parser, err
	}
}

func (p *Parser) shouldExclude(path string) bool {
	for _, exclusion := range p.exclude {
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
	// those as if they where normal imports.
	//
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
	exports, err := p.parseExports(n.Id, false, nil)
	if err != nil {
		return nil, err
	}
	n.AddErrors(exports.Errors...)
	for el := exports.Exports.Front(); el != nil; el = el.Next() {
		if el.Value != n.Id {
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

		// NOTE: at this point p.unwrapProxyExports is always true.
		exports, err = p.parseExports(importEntry.Path, p.unwrapProxyExports, nil)
		if err != nil {
			return nil, err
		}
		n.AddErrors(exports.Errors...)
		if importEntry.All {
			// If all imported, then dump every path in the resolved imports.
			for el := exports.Exports.Front(); el != nil; el = el.Next() {
				resolvedImports.Set(el.Value, true)
			}
			continue
		} else if len(importEntry.Names) == 0 {
			resolvedImports.Set(importEntry.Path, true)
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

	deps := make([]*graph.Node[*FileInfo], 0)
	for _, imported := range resolvedImports.Keys() {
		file, err := p.parseFile(imported)
		if err != nil {
			n.AddErrors(err)
			continue
		}
		node := graph.MakeNode(imported, file)
		if !p.shouldExclude(p.Display(node).Name) {
			deps = append(deps, node)
		}
	}
	return deps, nil
}

func (p *Parser) Display(n *graph.Node[*FileInfo]) graph.DisplayResult {
	return p.lang.Display(n.Id)
}
