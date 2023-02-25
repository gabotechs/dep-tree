package language

import (
	"context"
	"errors"
	"time"

	"dep-tree/internal/graph"
)

type TestLanguageData struct{}

type TestLanguage struct {
	imports map[string]*ImportsResult
	exports map[string]*ExportsResult
}

func (t *TestLanguage) MakeNode(id string) (*graph.Node[TestLanguageData], error) {
	return &graph.Node[TestLanguageData]{
		Id:     id,
		Errors: make([]error, 0),
		Data:   TestLanguageData{},
	}, nil
}

func (t *TestLanguage) ParseImports(ctx context.Context, id string) (context.Context, *ImportsResult, error) {
	time.Sleep(time.Millisecond)
	if imports, ok := t.imports[id]; ok {
		return ctx, imports, nil
	} else {
		return ctx, imports, errors.New(id + " not found")
	}
}

func (t *TestLanguage) ParseExports(ctx context.Context, id string) (context.Context, *ExportsResult, error) {
	time.Sleep(time.Millisecond)
	if exports, ok := t.exports[id]; ok {
		return ctx, exports, nil
	} else {
		return ctx, exports, errors.New(id + " not found")
	}
}
