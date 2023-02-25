package language

import "context"

type ExportsResult struct {
	// Exports: map from exported name to the absolute path from where it is exported
	//  NOTE: even though it could work returning a path relative to the file, it should return absolute
	Exports map[string]string
	// Errors: errors while parsing exports.
	Errors []error
}

type ExportsCacheKey string

func (p *Parser[T, F]) CachedParseExports(
	ctx context.Context,
	filePath string,
) (context.Context, *ExportsResult, error) {
	cacheKey := ExportsCacheKey(filePath)
	if cached, ok := ctx.Value(cacheKey).(*ExportsResult); ok {
		return ctx, cached, nil
	} else {
		ctx, file, err := p.CachedParseFile(ctx, filePath)
		if err != nil {
			return ctx, nil, err
		}
		ctx, result, err := p.lang.ParseExports(ctx, file)
		if err != nil {
			return ctx, nil, err
		}
		ctx = context.WithValue(ctx, cacheKey, result)
		return ctx, result, err
	}
}
