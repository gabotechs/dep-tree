package language

import (
	"errors"
	"fmt"

	"github.com/elliotchance/orderedmap/v2"

	"github.com/gabotechs/dep-tree/internal/utils"
)

func (en *ExportSymbol) name() string {
	if en.Alias != "" {
		return en.Alias
	} else {
		return en.Original
	}
}

// ExportEntries is the result of gathering all the export statements from
// a source file, in case the language implementation explicitly exports certain files.
type ExportEntries struct {
	// Symbols is an ordered map data structure where the keys are the symbols exported from
	// the source file and the values are path from where they are declared. Symbols might
	// be declared in a different path from where they are exported, for example:
	//
	// export { foo } from './bar'
	//
	// the `foo` symbol is being exported from the current file, but it's declared on the
	// `bar.ts` file.
	Symbols *orderedmap.OrderedMap[string, string]
	// Errors are the non-fatal errors that occurred while parsing exports. These
	// might be rendered nicely in a UI.
	Errors []error
}

func (p *Parser) parseExports(
	id string,
	unwrappedExports bool,
	stack *utils.CallStack,
) (*ExportEntries, error) {
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
		if export.AbsPath == id {
			for _, name := range export.Symbols {
				exports.Set(name.name(), export.AbsPath)
			}
			continue
		}

		var unwrapped *ExportEntries
		unwrapped, err = p.parseExports(export.AbsPath, unwrappedExports, stack)
		if err != nil {
			exportErrors = append(exportErrors, err)
			continue
		}

		if export.All {
			for el := unwrapped.Symbols.Front(); el != nil; el = el.Next() {
				if unwrappedExports {
					exports.Set(el.Key, el.Value)
				} else {
					exports.Set(el.Key, export.AbsPath)
				}
			}
			continue
		}
		exportErrors = append(exportErrors, unwrapped.Errors...)

		for _, name := range export.Symbols {
			if exportPath, ok := unwrapped.Symbols.Get(name.Original); ok {
				if unwrappedExports {
					exports.Set(name.name(), exportPath)
				} else {
					exports.Set(name.name(), export.AbsPath)
				}
			} else {
				exports.Set(name.name(), export.AbsPath)
				// errors = append(errors, fmt.Errorf(`name "%s" exported in "%s" from "%s" cannot be found in origin file`, name.Original, id, export.Id)).
			}
		}
	}

	result := ExportEntries{Symbols: exports, Errors: exportErrors}
	p.ExportsCache[cacheKey] = &result
	return &result, nil
}
