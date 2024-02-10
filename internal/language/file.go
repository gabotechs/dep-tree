package language

func (p *Parser) parseFile(absPath string) (*FileInfo, error) {
	if cached, ok := p.FileCache[absPath]; ok {
		return cached, nil
	}
	result, err := p.Lang.ParseFile(absPath)
	if err != nil {
		return nil, err
	}
	p.FileCache[absPath] = result
	return result, err
}
