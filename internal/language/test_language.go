package language

import (
	"errors"
	"time"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
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

func (t *TestLanguage) testParser() *Parser[TestFile] {
	return &Parser[TestFile]{
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

func (t *TestLanguage) Display(id string) dep_tree.DisplayResult {
	return dep_tree.DisplayResult{Name: id}
}
