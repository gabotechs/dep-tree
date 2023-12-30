package language

type FileCacheKey string

func (p *Parser[F]) parseFile(id string) (*F, error) {
	if cached, ok := p.fileCache[id]; ok {
		return cached, nil
	}
	result, err := p.lang.ParseFile(id)
	if err != nil {
		return nil, err
	}
	p.fileCache[id] = result
	return result, err
}
