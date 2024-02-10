package language

import (
	"errors"
	"fmt"

	"github.com/elliotchance/orderedmap/v2"

	"github.com/gabotechs/dep-tree/internal/utils"
)

type ExportName struct {
	Original string
	Alias    string
}

func (en *ExportName) name() string {
	if en.Alias != "" {
		return en.Alias
	} else {
		return en.Original
	}
}

type ExportEntry struct {
	// All: all the names from Path are exported.
	All bool
	// Names: exported specific names from Path.
	Names []ExportName
	// Path: absolute path from where they are exported, it might be from the same file or from another.
	Path string
}

type ExportsEntries struct {
	// Exports: array of ExportEntry
	//  NOTE: even though it could work returning a path relative to the file, it should return absolute.
	Exports []ExportEntry
	// Errors: errors while parsing exports.
	Errors []error
}

func (p *Parser) parseExports(
	id string,
	unwrappedExports bool,
	stack *utils.CallStack,
) (*ExportsResult, error) {
	if stack == nil {
		stack = utils.NewCallStack()
	}
	if err := stack.Push(id); err != nil {
		return nil, errors.New("circular export: " + err.Error())
	}
	defer stack.Pop()
	cacheKey := fmt.Sprintf("%s-%t", id, unwrappedExports)
	if cached, ok := p.ExportsCache[cacheKey]; ok {
		return cached, nil
	}

	file, err := p.parseFile(id)
	if err != nil {
		return nil, err
	}

	wrapped, err := p.Lang.ParseExports(file)
	if err != nil {
		return nil, err
	}

	exports := orderedmap.NewOrderedMap[string, string]()
	var exportErrors []error

	for _, export := range wrapped.Exports {
		if export.Path == id {
			for _, name := range export.Names {
				exports.Set(name.name(), export.Path)
			}
			continue
		}

		var unwrapped *ExportsResult
		unwrapped, err = p.parseExports(export.Path, unwrappedExports, stack)
		if err != nil {
			exportErrors = append(exportErrors, err)
			continue
		}

		if export.All {
			for el := unwrapped.Symbols.Front(); el != nil; el = el.Next() {
				if unwrappedExports {
					exports.Set(el.Key, el.Value)
				} else {
					exports.Set(el.Key, export.Path)
				}
			}
			continue
		}
		exportErrors = append(exportErrors, unwrapped.Errors...)

		for _, name := range export.Names {
			if exportPath, ok := unwrapped.Symbols.Get(name.Original); ok {
				if unwrappedExports {
					exports.Set(name.name(), exportPath)
				} else {
					exports.Set(name.name(), export.Path)
				}
			} else {
				exports.Set(name.name(), export.Path)
				// errors = append(errors, fmt.Errorf(`name "%s" exported in "%s" from "%s" cannot be found in origin file`, name.Original, id, export.Id)).
			}
		}
	}

	result := ExportsResult{Symbols: exports, Errors: exportErrors}
	p.ExportsCache[cacheKey] = &result
	return &result, nil
}
