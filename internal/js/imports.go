package js

import (
	"context"
	"path"

	"github.com/elliotchance/orderedmap/v2"

	"dep-tree/internal/js/grammar"
)

type ImportsCacheKey string

func (p *Parser) parseImports(
	ctx context.Context,
	filePath string,
) (context.Context, *orderedmap.OrderedMap[string, []string], error) {
	cacheKey := ImportsCacheKey(filePath)
	if cached, ok := ctx.Value(cacheKey).(*orderedmap.OrderedMap[string, []string]); ok {
		return ctx, cached, nil
	} else {
		ctx, result, err := p.uncachedParseImports(ctx, filePath)
		if err != nil {
			return ctx, nil, err
		}
		ctx = context.WithValue(ctx, cacheKey, result)
		return ctx, result, err
	}
}

func (p *Parser) uncachedParseImports(
	ctx context.Context,
	filePath string,
) (context.Context, *orderedmap.OrderedMap[string, []string], error) {
	ctx, jsFile, err := grammar.Parse(ctx, filePath)
	if err != nil {
		return ctx, nil, err
	}

	imports := orderedmap.NewOrderedMap[string, []string]()

	for _, stmt := range jsFile.Statements {
		var importPath string
		var names []string
		switch {
		case stmt == nil:
			continue
		case stmt.StaticImport != nil:
			importPath = stmt.StaticImport.Path
			names = gatherNamesFromStaticImport(stmt.StaticImport)
		case stmt.DynamicImport != nil:
			importPath = stmt.DynamicImport.Path
			names = []string{"*"}
		default:
			continue
		}
		var resolvedPath string
		ctx, resolvedPath, err = p.ResolvePath(ctx, importPath, path.Dir(filePath))
		if err != nil {
			return ctx, nil, err
		} else if resolvedPath != "" {
			imports.Set(resolvedPath, names)
		}
	}
	return ctx, imports, nil
}

func gatherNamesFromStaticImport(si *grammar.StaticImport) []string {
	names := make([]string, 0)

	if imported := si.Imported; imported != nil {
		if imported.Default {
			names = append(names, "default")
		}
		if selection := imported.SelectionImport; selection != nil {
			if selection.AllImport != nil {
				names = append(names, "*")
			}
			if selection.Deconstruction != nil {
				names = append(names, selection.Deconstruction.Names...)
			}
		}
	} else {
		names = append(names, "*")
	}
	return names
}
