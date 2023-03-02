package language

import (
	"context"
	"fmt"
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
	result, err := p.lang.ParseExports(file)
	if err != nil {
		return ctx, nil, err
	}
	ctx = context.WithValue(ctx, cacheKey, result)
	return ctx, result, err
}

type UnwrappedExportsResult struct {
	// Exports: map from exported name to exported path.
	Exports map[string]string
	// Errors: errors gathered while resolving exports.
	Errors []error
}

type UnwrappedExportsCacheKey string

func (p *Parser[T, F]) CachedUnwrappedParseExports(
	ctx context.Context,
	id string,
) (context.Context, *UnwrappedExportsResult, error) {
	unwrappedCacheKey := UnwrappedExportsCacheKey(id)
	if cached, ok := ctx.Value(unwrappedCacheKey).(*UnwrappedExportsResult); ok {
		return ctx, cached, nil
	}
	ctx, wrapped, err := p.CachedParseExports(ctx, id)
	if err != nil {
		return ctx, nil, err
	}
	exports := make(map[string]string)
	var errors []error

	for _, export := range wrapped.Exports {
		if export.Id == id {
			for _, name := range export.Names {
				exports[name.name()] = export.Id
			}
			continue
		}

		var unwrapped *UnwrappedExportsResult
		ctx, unwrapped, err = p.CachedUnwrappedParseExports(ctx, export.Id)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		if export.All {
			for exportedName, exportId := range unwrapped.Exports {
				exports[exportedName] = exportId
			}
			continue
		}

		for _, name := range export.Names {
			if exportId, ok := unwrapped.Exports[name.Original]; ok {
				exports[name.name()] = exportId
			} else {
				errors = append(errors, fmt.Errorf(`name "%s" exported in "%s" from "%s" cannot be found in origin file`, name.Original, id, export.Id))
			}
		}
	}

	return ctx, &UnwrappedExportsResult{Exports: exports, Errors: errors}, nil
}
