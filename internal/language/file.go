package language

import "context"

type FileCacheKey string

func (p *Parser[F]) parseFile(ctx context.Context, id string) (context.Context, *F, error) {
	cacheKey := FileCacheKey(id)
	if cached, ok := ctx.Value(cacheKey).(*F); ok {
		return ctx, cached, nil
	}
	result, err := p.lang.ParseFile(id)
	if err != nil {
		return ctx, nil, err
	}
	return context.WithValue(ctx, cacheKey, result), result, err
}
