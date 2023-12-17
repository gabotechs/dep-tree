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

type ImportsResult struct {
	// Imports: ordered map from absolute imported path to the array of names that where imported.
	//  if one of the names is *, then all the names are imported
	Imports []ImportEntry
	// Errors: errors while parsing imports.
	Errors []error
}

type ImportsCacheKey string

func (p *Parser[F]) CachedParseImports(
	ctx context.Context,
	id string,
) (context.Context, *ImportsResult, error) {
	cacheKey := ImportsCacheKey(id)
	if cached, ok := ctx.Value(cacheKey).(*ImportsResult); ok {
		return ctx, cached, nil
	} else {
		ctx, file, err := p.CachedParseFile(ctx, id)
		if err != nil {
			return ctx, nil, err
		}
		result, err := p.lang.ParseImports(file)
		if err != nil {
			return ctx, nil, err
		}
		ctx = context.WithValue(ctx, cacheKey, result)
		return ctx, result, err
	}
}
