package language

import (
	"context"

	"github.com/elliotchance/orderedmap/v2"
)

type ImportEntry struct {
	All   bool
	Names []string
}

type ImportsResult struct {
	// Imports: ordered map from absolute imported path to the array of names that where imported.
	//  if one of the names is *, then all the names are imported
	Imports *orderedmap.OrderedMap[string, ImportEntry]
	// Errors: errors while parsing imports.
	Errors []error
}

type ImportsCacheKey string

func (p *Parser[T, F]) CachedParseImports(
	ctx context.Context,
	filePath string,
) (context.Context, *ImportsResult, error) {
	cacheKey := ImportsCacheKey(filePath)
	if cached, ok := ctx.Value(cacheKey).(*ImportsResult); ok {
		return ctx, cached, nil
	} else {
		ctx, file, err := p.CachedParseFile(ctx, filePath)
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
