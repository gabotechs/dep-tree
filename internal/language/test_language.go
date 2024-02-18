package language

import (
	"errors"
	"time"
)

type TestFileContent struct {
	Name string
}

type TestLanguage struct {
	imports map[string]*ImportsResult
	exports map[string]*ExportsResult
}

func (t *TestLanguage) testParser() *Parser {
	return &Parser{
		Lang:         t,
		FileCache:    map[string]*FileInfo{},
		ImportsCache: map[string]*ImportsResult{},
		ExportsCache: map[string]*ExportEntries{},
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

func (t *TestLanguage) ParseExports(file *FileInfo) (*ExportsResult, error) {
	time.Sleep(time.Millisecond)
	content := file.Content.(TestFileContent)
	if exports, ok := t.exports[content.Name]; ok {
		return exports, nil
	} else {
		return exports, errors.New(content.Name + " not found")
	}
}
