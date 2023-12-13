package language

import (
	"context"
	"errors"
	"time"

	"github.com/gabotechs/dep-tree/internal/graph"
)

type TestFile struct {
	Name string
}

type TestLanguageData struct{}

type TestLanguage struct {
	imports map[string]*ImportsResult
	exports map[string]*ExportsResult
}

func (t *TestLanguage) testParser(entrypoint string) *Parser[TestLanguageData, TestFile] {
	_, parser, _ := makeParser(context.Background(), entrypoint, func(ctx context.Context, _ string, _ *struct{}) (context.Context, Language[TestLanguageData, TestFile], error) {
		return ctx, t, nil
	}, nil)
	return parser
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

func (t *TestLanguage) ParseImports(file *TestFile) (*ImportsResult, error) {
	time.Sleep(time.Millisecond)
	if imports, ok := t.imports[file.Name]; ok {
		return imports, nil
	} else {
		return imports, errors.New(file.Name + " not found")
	}
}

func (t *TestLanguage) ParseExports(file *TestFile) (*ExportsResult, error) {
	time.Sleep(time.Millisecond)
	if exports, ok := t.exports[file.Name]; ok {
		return exports, nil
	} else {
		return exports, errors.New(file.Name + " not found")
	}
}
