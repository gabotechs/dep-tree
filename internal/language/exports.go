package language

import (
	"context"
	"fmt"

	"github.com/elliotchance/orderedmap/v2"

	"dep-tree/internal/utils"
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
	// Id: absolute path from where they are exported, it might be from the same file or from another.
	Id string
}

type ExportsResult struct {
	// Exports: array of ExportEntry
	//  NOTE: even though it could work returning a path relative to the file, it should return absolute.
	Exports []ExportEntry
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
	}
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

type UnwrappedExportsResult struct {
	// Exports: map from exported name to exported path.
	Exports *orderedmap.OrderedMap[string, string]
	// Errors: errors gathered while resolving exports.
	Errors []error
}

type UnwrappedExportsCacheKey string

func (p *Parser[T, F]) CachedUnwrappedParseExports(
	ctx context.Context,
	id string,
) (context.Context, *UnwrappedExportsResult, error) {
	return p.cachedUnwrappedParseExports(ctx, id, make(map[string]bool))
}

func (p *Parser[T, F]) cachedUnwrappedParseExports(
	ctx context.Context,
	id string,
	seen map[string]bool,
) (context.Context, *UnwrappedExportsResult, error) {
	if _, ok := seen[id]; ok {
		return ctx, nil, fmt.Errorf("circular export starting and ending on %s", id)
	} else {
		seenCopy := utils.Merge(nil, seen)
		seenCopy[id] = true
		seen = seenCopy
	}

	unwrappedCacheKey := UnwrappedExportsCacheKey(id)
	if cached, ok := ctx.Value(unwrappedCacheKey).(*UnwrappedExportsResult); ok {
		return ctx, cached, nil
	}
	ctx, wrapped, err := p.CachedParseExports(ctx, id)
	if err != nil {
		return ctx, nil, err
	}
	exports := orderedmap.NewOrderedMap[string, string]()
	var errors []error

	for _, export := range wrapped.Exports {
		if export.Id == id {
			for _, name := range export.Names {
				exports.Set(name.name(), export.Id)
			}
			continue
		}

		var unwrapped *UnwrappedExportsResult
		ctx, unwrapped, err = p.cachedUnwrappedParseExports(ctx, export.Id, seen)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		if export.All {
			exports = unwrapped.Exports
			continue
		}
		errors = append(errors, unwrapped.Errors...)

		for _, name := range export.Names {
			if exportId, ok := unwrapped.Exports.Get(name.Original); ok {
				exports.Set(name.name(), exportId)
			} else {
				errors = append(errors, fmt.Errorf(`name "%s" exported in "%s" from "%s" cannot be found in origin file`, name.Original, id, export.Id))
			}
		}
	}

	return ctx, &UnwrappedExportsResult{Exports: exports, Errors: errors}, nil
}
