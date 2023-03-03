package language

import (
	"context"
	"errors"
	"time"

	"dep-tree/internal/graph"
)

type TestFile struct {
	Name string
}

type TestLanguageData struct{}

type TestLanguage struct {
	imports map[string]*ImportsResult
	exports map[string]*ExportsResult
}

var _ Language[TestLanguageData, TestFile] = &TestLanguage{}

func (t *TestLanguage) ParseFile(id string) (*TestFile, error) {
	time.Sleep(time.Millisecond)
	return &TestFile{
		Name: id,
	}, nil
}

func (t *TestLanguage) MakeNode(id string) (*graph.Node[TestLanguageData], error) {
	return &graph.Node[TestLanguageData]{
		Id:     id,
		Errors: make([]error, 0),
		Data:   TestLanguageData{},
	}, nil
}

func (t *TestLanguage) ParseImports(ctx context.Context, file *TestFile) (context.Context, *ImportsResult, error) {
	time.Sleep(time.Millisecond)
	if imports, ok := t.imports[file.Name]; ok {
		return ctx, imports, nil
	} else {
		return ctx, imports, errors.New(file.Name + " not found")
	}
}

func (t *TestLanguage) ParseExports(ctx context.Context, file *TestFile) (context.Context, *ExportsResult, error) {
	time.Sleep(time.Millisecond)
	if exports, ok := t.exports[file.Name]; ok {
		return ctx, exports, nil
	} else {
		return ctx, exports, errors.New(file.Name + " not found")
	}
}
