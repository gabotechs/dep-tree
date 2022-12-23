package js

import (
	"context"
	"fmt"
	"path"

	"dep-tree/internal/js/grammar"
	"dep-tree/internal/utils"
)

type ExportsCacheKey string

func (p *Parser) parseExports(
	ctx context.Context,
	filePath string,
) (context.Context, map[string]string, error) {
	cacheKey := ExportsCacheKey(filePath)
	if cached, ok := ctx.Value(cacheKey).(map[string]string); ok {
		return ctx, cached, nil
	} else {
		ctx, result, err := p.uncachedParseExports(ctx, filePath)
		if err != nil {
			return ctx, nil, err
		}
		ctx = context.WithValue(ctx, cacheKey, result)
		return ctx, result, err
	}
}

func (p *Parser) uncachedParseExports(
	ctx context.Context,
	filePath string,
) (context.Context, map[string]string, error) {
	ctx, jsFile, err := grammar.Parse(ctx, filePath)
	if err != nil {
		return ctx, nil, err
	}
	exported := make(map[string]string)
	for _, stmt := range jsFile.Statements {
		switch {
		case stmt == nil:
			continue
		case stmt.DeclarationExport != nil:
			exported[stmt.DeclarationExport.Name] = filePath
		case stmt.ListExport != nil:
			if stmt.ListExport.ExportDeconstruction != nil {
				for _, name := range stmt.ListExport.ExportDeconstruction.Names {
					exportedName := name.Alias
					if exportedName == "" {
						exportedName = name.Original
					}
					exported[exportedName] = filePath
				}
			}
		case stmt.DefaultExport != nil:
			if stmt.DefaultExport.Default {
				exported["default"] = filePath
			}
		case stmt.ProxyExport != nil:
			var exportFrom string
			ctx, exportFrom, err = p.ResolvePath(ctx, stmt.ProxyExport.From, path.Dir(filePath))
			if err != nil {
				return ctx, nil, err
			}
			var proxyExports map[string]string
			ctx, proxyExports, err = p.parseExports(ctx, exportFrom)
			if err != nil {
				return ctx, nil, err
			}
			if stmt.ProxyExport.ExportAll {
				if stmt.ProxyExport.ExportAllAlias != "" {
					exported[stmt.ProxyExport.ExportAllAlias] = filePath
				} else {
					exported = utils.Merge(exported, proxyExports)
				}
			} else if stmt.ProxyExport.ExportDeconstruction != nil {
				for _, name := range stmt.ProxyExport.ExportDeconstruction.Names {
					alias := name.Alias
					original := name.Original
					if alias == "" {
						alias = original
					}
					if proxyPath, ok := proxyExports[original]; ok {
						exported[alias] = proxyPath
					} else {
						return ctx, nil, fmt.Errorf("cannot import \"%s\" from %s", original, exportFrom)
					}
				}
			}
		default:
			continue
		}
	}
	return ctx, exported, nil
}
