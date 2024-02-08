package language

type FileCacheKey string

func (p *Parser) parseFile(absPath string) (*FileInfo, error) {
	if cached, ok := p.fileCache[absPath]; ok {
		return cached, nil
	}
	result, err := p.lang.ParseFile(absPath)
	if err != nil {
		return nil, err
	}
	p.fileCache[absPath] = result
	return result, err
}
