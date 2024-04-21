package dart

import (
	"bytes"
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
			var importPath string

			if statement.Import.IsAbsolute {
				// Code files must always be in the <root-where-pubspec-is-located>/lib directory.
				importPath = filepath.Join(findClosestDartRootDir(file.AbsPath), "lib", statement.Import.From)
			} else {
				// Relative imports are relative to the current file.
				importPath = filepath.Join(filepath.Dir(file.AbsPath), statement.Import.From)
			}

			// fmt.Println(importPath)
			result.Imports = append(result.Imports, language.ImportEntry{
				AbsPath: importPath,
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
