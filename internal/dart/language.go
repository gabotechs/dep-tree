package dart

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/language"
)

var Extensions = []string{"dart"}

type Language struct{}

func (l *Language) ParseFile(path string) (*language.FileInfo, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	file, err := ParseFile(path)
	if err != nil {
		return nil, err
	}
	currentDir, _ := os.Getwd()
	relPath, _ := filepath.Rel(currentDir, path)
	return &language.FileInfo{
		Content: file.Statements,                    // dump the parsed statements into the FileInfo struct.
		Loc:     bytes.Count(content, []byte("\n")), // get the amount of lines of code.
		Size:    len(content),                       // get the size of the file in bytes.
		AbsPath: path,                               // provide its absolute path.
		RelPath: relPath,                            // provide the path relative to the current dir.
	}, nil
}

func (l *Language) ParseImports(file *language.FileInfo) (*language.ImportsResult, error) {
	var result language.ImportsResult

	for _, statement := range file.Content.([]Statement) {
		if statement.Import != nil {
			fmt.Println(filepath.Join(filepath.Dir(file.AbsPath), statement.Import.From))
			result.Imports = append(result.Imports, language.ImportEntry{
				AbsPath: filepath.Join(filepath.Dir(file.AbsPath), statement.Import.From),
			})
		}
	}

	return &result, nil
}

func (l *Language) ParseExports(file *language.FileInfo) (*language.ExportsResult, error) {
	var result language.ExportsResult

	for _, statement := range file.Content.([]Statement) {
		if statement.Export != nil {
			result.Exports = append(result.Exports, language.ExportEntry{
				AbsPath: filepath.Join(filepath.Dir(file.AbsPath), statement.Export.From),
			})
		}
	}

	return &result, nil
}
