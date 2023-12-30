package python

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python/python_grammar"
	"github.com/gabotechs/dep-tree/internal/utils"
)

var Extensions = []string{
	"py",
	"pyi",
	"pyx",
}

type Language struct {
	PythonPath []string
	cfg        *Config
}

var _ language.Language[python_grammar.File] = &Language{}

var rootFiles = []string{
	"pyproject.toml",
	"setup.py",
	"poetry.toml",
	"poetry.lock",
	"requirements.txt",
	".pylintrc",
	".git/index",
}

func isRootFilePresent(dir string) bool {
	for _, rootFile := range rootFiles {
		if utils.FileExists(path.Join(dir, rootFile)) {
			return true
		}
	}
	return false
}

func MakePythonLanguage(entrypoint string, cfg *Config) (language.Language[python_grammar.File], error) {
	lang := Language{
		cfg: cfg,
	}
	if lang.cfg == nil {
		lang.cfg = &Config{}
	}
	entrypointAbsPath, err := filepath.Abs(entrypoint)
	if err != nil {
		return nil, err
	}
	var baseDir string
	switch {
	case utils.FileExists(entrypointAbsPath):
		baseDir = filepath.Dir(entrypointAbsPath)
	case utils.DirExists(entrypointAbsPath):
		baseDir = entrypointAbsPath
	default:
		return nil, fmt.Errorf("file %s does not exist", entrypoint)
	}
	lookupDir := baseDir
	rootFilePresent := isRootFilePresent(lookupDir)
	for !rootFilePresent && len(lookupDir) > 2 {
		lookupDir = path.Dir(lookupDir)
		rootFilePresent = isRootFilePresent(lookupDir)
	}
	// Search for the root path based on some key files.
	if rootFilePresent {
		lang.PythonPath = append(lang.PythonPath, lookupDir)
	}

	// Add the entrypoint's directory.
	lang.PythonPath = append(lang.PythonPath, baseDir)

	// Add anything present on the PYTHONPATH.
	pp := os.Getenv("PYTHONPATH")
	if pp != "" {
		lang.PythonPath = append(lang.PythonPath, strings.Split(pp, ":")...)
	}
	return &lang, nil
}

func (l *Language) ParseFile(id string) (*python_grammar.File, error) {
	return python_grammar.Parse(id)
}
