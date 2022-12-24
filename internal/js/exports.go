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
			// Is this even possible?
		case stmt.DeclarationExport != nil:
			handleDeclarationExport(stmt.DeclarationExport, filePath, exported)
		case stmt.ListExport != nil:
			handleListExport(stmt.ListExport, filePath, exported)
		case stmt.DefaultExport != nil:
			handleDefaultExport(stmt.DefaultExport, filePath, exported)
		case stmt.ProxyExport != nil:
			ctx, err = p.handleProxyExport(ctx, stmt.ProxyExport, filePath, exported)
		}
		if err != nil {
			return ctx, nil, err
		}
	}
	return ctx, exported, nil
}

func handleDeclarationExport(
	stmt *grammar.DeclarationExport,
	filePath string,
	dumpOn map[string]string,
) {
	dumpOn[stmt.Name] = filePath
}

func handleListExport(
	stmt *grammar.ListExport,
	filePath string,
	dumpOn map[string]string,
) {
	if stmt.ExportDeconstruction != nil {
		for _, name := range stmt.ExportDeconstruction.Names {
			exportedName := name.Alias
			if exportedName == "" {
				exportedName = name.Original
			}
			dumpOn[exportedName] = filePath
		}
	}
}

func handleDefaultExport(
	stmt *grammar.DefaultExport,
	filePath string,
	dumpOn map[string]string,
) {
	if stmt.Default {
		dumpOn["default"] = filePath
	}
}

func (p *Parser) handleProxyExport(
	ctx context.Context,
	stmt *grammar.ProxyExport,
	filePath string,
	dumpOn map[string]string,
) (context.Context, error) {
	ctx, exportFrom, err := p.ResolvePath(ctx, stmt.From, path.Dir(filePath))
	if err != nil {
		return ctx, err
	}
	// WARN: this call is recursive, be aware!!!
	ctx, proxyExports, err := p.parseExports(ctx, exportFrom)
	switch {
	case err != nil:
		return ctx, err
	case stmt.ExportAll:
		if stmt.ExportAllAlias != "" {
			dumpOn[stmt.ExportAllAlias] = filePath
		} else {
			utils.Merge(dumpOn, proxyExports)
		}
	case stmt.ExportDeconstruction != nil:
		for _, name := range stmt.ExportDeconstruction.Names {
			if proxyPath, ok := proxyExports[name.Original]; ok {
				dumpOn[name.AliasOrOriginal()] = proxyPath
			} else {
				return ctx, fmt.Errorf("cannot import \"%s\" from %s", name.Original, exportFrom)
			}
		}
	}
	return ctx, nil
}
