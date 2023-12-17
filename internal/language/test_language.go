package language

import (
	"context"
	"errors"
	"time"
)

type TestFile struct {
	Name string
}

func (t TestFile) Loc() int {
	return 0
}

func (t TestFile) Size() int {
	return 0
}

type TestLanguage struct {
	imports map[string]*ImportsResult
	exports map[string]*ExportsResult
}

func (t *TestLanguage) testParser(entrypoint string) *Parser[TestFile] {
	_, parser, _ := makeParser(context.Background(), entrypoint, func(ctx context.Context, _ string, _ *struct{}) (context.Context, Language[TestFile], error) {
		return ctx, t, nil
	}, nil)
	return parser
}

var _ Language[TestFile] = &TestLanguage{}

func (t *TestLanguage) ParseFile(id string) (*TestFile, error) {
	time.Sleep(time.Millisecond)
	return &TestFile{
		Name: id,
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
