package language

import (
	"context"
)

type ImportEntry struct {
	// All: if all the names from Path are imported.
	All bool
	// Names: what specific names form Path are imported.
	Names []string
	// Path: from where are the names imported.
	Path string
}

func AllImport(path string) ImportEntry {
	return ImportEntry{All: true, Path: path}
}

func EmptyImport(path string) ImportEntry {
	return ImportEntry{Path: path}
}

func NamesImport(names []string, path string) ImportEntry {
	return ImportEntry{Names: names, Path: path}
}

type ImportsResult struct {
	// Imports: ordered map from absolute imported path to the array of names that where imported.
	//  if one of the names is *, then all the names are imported
	Imports []ImportEntry
	// Errors: errors while parsing imports.
	Errors []error
}

type ImportsCacheKey string

func (p *Parser[F]) gatherImportsFromFile(
	ctx context.Context,
	id string,
) (context.Context, *ImportsResult, error) {
	cacheKey := ImportsCacheKey(id)
	if cached, ok := ctx.Value(cacheKey).(*ImportsResult); ok {
		return ctx, cached, nil
	}
	ctx, file, err := p.parseFile(ctx, id)
	if err != nil {
		return ctx, nil, err
	}
	result, err := p.lang.ParseImports(file)
	if err != nil {
		return ctx, nil, err
	}
	return context.WithValue(ctx, cacheKey, result), result, err
}
