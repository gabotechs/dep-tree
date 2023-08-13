package python

import (
	"dep-tree/internal/utils"
	"os"
	"path/filepath"
	"strings"

	"dep-tree/internal/language"
	"dep-tree/internal/python/python_grammar"
)

type Language struct {
	PythonPath []string
}

var _ language.Language[Data, python_grammar.File] = &Language{}

func MakePythonLanguage(entrypoint string) (language.Language[Data, python_grammar.File], error) {
	entrypointAbsPath, err := filepath.Abs(entrypoint)
	if err != nil {
		return nil, err
	}
	var baseDir string
	if utils.FileExists(entrypointAbsPath) {
		baseDir = filepath.Dir(entrypointAbsPath)
	} else if utils.DirExists(entrypointAbsPath) {
		baseDir = entrypointAbsPath
	}

	pp := os.Getenv("PYTHONPATH")
	lang := Language{
		PythonPath: []string{baseDir},
	}
	if pp != "" {
		lang.PythonPath = strings.Split(pp, ":")
	}
	return &lang, nil
}

func (l *Language) ParseFile(id string) (*python_grammar.File, error) {
	return python_grammar.Parse(id)
}
