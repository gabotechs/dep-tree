package language

import "context"

type FileCacheKey string

func (p *Parser[T, F]) CachedParseFile(ctx context.Context, id string) (context.Context, *F, error) {
	cacheKey := FileCacheKey(id)
	if cached, ok := ctx.Value(cacheKey).(*F); ok {
		return ctx, cached, nil
	} else {
		result, err := p.lang.ParseFile(id)
		if err != nil {
			return ctx, nil, err
		}
		ctx = context.WithValue(ctx, cacheKey, result)
		return ctx, result, err
	}
}
