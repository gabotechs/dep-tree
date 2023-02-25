package js

import (
	"context"
	"fmt"
	"path"

	"dep-tree/internal/js/grammar"
	"dep-tree/internal/language"
	"dep-tree/internal/utils"
)

type ExportsCacheKey string

func (l *Language) ParseExports(
	ctx context.Context,
	file *grammar.File,
) (context.Context, *language.ExportsResult, error) {
	results := &language.ExportsResult{
		Exports: make(map[string]string),
		Errors:  []error{},
	}

	for _, stmt := range file.Statements {
		switch {
		case stmt == nil:
			// Is this even possible?
		case stmt.DeclarationExport != nil:
			handleDeclarationExport(stmt.DeclarationExport, file.Path, results.Exports)
		case stmt.ListExport != nil:
			handleListExport(stmt.ListExport, file.Path, results.Exports)
		case stmt.DefaultExport != nil:
			handleDefaultExport(stmt.DefaultExport, file.Path, results.Exports)
		case stmt.ProxyExport != nil:
			var err error
			ctx, err = l.handleProxyExport(ctx, stmt.ProxyExport, file.Path, results.Exports)
			if err != nil {
				results.Errors = append(results.Errors, err)
			}
		}
	}
	return ctx, results, nil
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

func (l *Language) handleProxyExport(
	ctx context.Context,
	stmt *grammar.ProxyExport,
	filePath string,
	dumpOn map[string]string,
) (context.Context, error) {
	ctx, exportFrom, err := l.ResolvePath(ctx, stmt.From, path.Dir(filePath))
	if err != nil {
		return ctx, err
	}
	// TODO: this should go away.
	parsed, err := l.ParseFile(exportFrom)
	if err != nil {
		return ctx, err
	}
	// WARN: this call is recursive, be aware!!!
	ctx, proxyExports, err := l.ParseExports(ctx, parsed)
	switch {
	case err != nil:
		return ctx, err
	case stmt.ExportAll:
		if stmt.ExportAllAlias != "" {
			dumpOn[stmt.ExportAllAlias] = filePath
		} else {
			utils.Merge(dumpOn, proxyExports.Exports)
		}
	case stmt.ExportDeconstruction != nil:
		for _, name := range stmt.ExportDeconstruction.Names {
			if proxyPath, ok := proxyExports.Exports[name.Original]; ok {
				dumpOn[name.AliasOrOriginal()] = proxyPath
			} else {
				return ctx, fmt.Errorf("cannot import \"%s\" from %s", name.Original, exportFrom)
			}
		}
	}
	return ctx, nil
}
