package language

import (
	"errors"
	"time"

	"github.com/gabotechs/dep-tree/internal/graph"
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
	exports map[string]*ExportsEntries
}

func (t *TestLanguage) testParser(entrypoint string) *Parser[TestFile] {
	entrypointNode := graph.MakeNode(entrypoint, FileInfo{})
	return &Parser[TestFile]{
		entrypoint:   entrypointNode,
		lang:         t,
		fileCache:    map[string]*TestFile{},
		importsCache: map[string]*ImportsResult{},
		exportsCache: map[string]*ExportsResult{},
	}
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

func (t *TestLanguage) ParseExports(file *TestFile) (*ExportsEntries, error) {
	time.Sleep(time.Millisecond)
	if exports, ok := t.exports[file.Name]; ok {
		return exports, nil
	} else {
		return exports, errors.New(file.Name + " not found")
	}
}
