package language

import (
	"errors"
	"time"

	"github.com/gabotechs/dep-tree/internal/graph"
)

type TestFileContent struct {
	Name string
}

type TestLanguage struct {
	imports map[string]*ImportsResult
	exports map[string]*ExportsEntries
}

func (t *TestLanguage) testParser() *Parser {
	return &Parser{
		lang:         t,
		fileCache:    map[string]*FileInfo{},
		importsCache: map[string]*ImportsResult{},
		exportsCache: map[string]*ExportsResult{},
	}
}

var _ Language = &TestLanguage{}

func (t *TestLanguage) ParseFile(id string) (*FileInfo, error) {
	time.Sleep(time.Millisecond)
	return &FileInfo{
		Content: TestFileContent{id},
	}, nil
}

func (t *TestLanguage) ParseImports(file *FileInfo) (*ImportsResult, error) {
	time.Sleep(time.Millisecond)
	content := file.Content.(TestFileContent)
	if imports, ok := t.imports[content.Name]; ok {
		return imports, nil
	} else {
		return imports, errors.New(content.Name + " not found")
	}
}

func (t *TestLanguage) ParseExports(file *FileInfo) (*ExportsEntries, error) {
	time.Sleep(time.Millisecond)
	content := file.Content.(TestFileContent)
	if exports, ok := t.exports[content.Name]; ok {
		return exports, nil
	} else {
		return exports, errors.New(content.Name + " not found")
	}
}

func (t *TestLanguage) Display(id string) graph.DisplayResult {
	return graph.DisplayResult{Name: id}
}
