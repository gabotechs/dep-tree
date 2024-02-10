package language

func (p *Parser) gatherImportsFromFile(id string) (*ImportsResult, error) {
	if cached, ok := p.ImportsCache[id]; ok {
		return cached, nil
	}
	file, err := p.parseFile(id)
	if err != nil {
		return nil, err
	}
	result, err := p.Lang.ParseImports(file)
	if err != nil {
		return nil, err
	}
	p.ImportsCache[id] = result
	return result, err
}
