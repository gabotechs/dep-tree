package language

import (
	"context"
)

type ImportEntry struct {
	// All: if all the names from Id are imported.
	All bool
	// Names: what specific names form Id are imported.
	Names []string
	// Id: from where are the names imported.
	Id string
}

type ImportsResult struct {
	// Imports: ordered map from absolute imported path to the array of names that where imported.
	//  if one of the names is *, then all the names are imported
	Imports []ImportEntry
	// Errors: errors while parsing imports.
	Errors []error
}

type ImportsCacheKey string

func (p *Parser[T, F]) CachedParseImports(
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
		ctx, result, err := p.lang.ParseImports(ctx, file)
		if err != nil {
			return ctx, nil, err
		}
		ctx = context.WithValue(ctx, cacheKey, result)
		return ctx, result, err
	}
}
