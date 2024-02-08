package language

type ImportEntry struct {
	// All: if all the names from Path are imported.
	All bool
	// Names: what specific names form Path are imported.
	Names []string
	// Path: from where are the names imported.
	Path string
}

func AllImport(path string) ImportEntry {
	return ImportEntry{All: true, Path: path}
}

func EmptyImport(path string) ImportEntry {
	return ImportEntry{Path: path}
}

func NamesImport(names []string, path string) ImportEntry {
	return ImportEntry{Names: names, Path: path}
}

type ImportsResult struct {
	// Imports: ordered map from absolute imported path to the array of names that where imported.
	//  if one of the names is *, then all the names are imported
	Imports []ImportEntry
	// Errors: errors while parsing imports.
	Errors []error
}

type ImportsCacheKey string

func (p *Parser) gatherImportsFromFile(id string) (*ImportsResult, error) {
	if cached, ok := p.importsCache[id]; ok {
		return cached, nil
	}
	file, err := p.parseFile(id)
	if err != nil {
		return nil, err
	}
	result, err := p.lang.ParseImports(file)
	if err != nil {
		return nil, err
	}
	p.importsCache[id] = result
	return result, err
}
