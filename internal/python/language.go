package python

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python/python_grammar"
)

var Extensions = []string{
	"py",
	"pyi",
	"pyx",
}

type Language struct {
	cfg *Config
}

var _ language.Language = &Language{}

func (l *Language) Display(id string) graph.DisplayResult {
	basePath := findClosestDirWithRootFile(filepath.Dir(id))
	result, err := filepath.Rel(basePath, id)
	if err != nil {
		return graph.DisplayResult{Name: id}
	}
	return graph.DisplayResult{Name: result}
}

func MakePythonLanguage(cfg *Config) (language.Language, error) {
	lang := Language{
		cfg: cfg,
	}
	if lang.cfg == nil {
		lang.cfg = &Config{}
	}

	// Add anything present on the PYTHONPATH.
	pp := os.Getenv("PYTHONPATH")
	if pp != "" {
		lang.cfg.PythonPath = append(lang.cfg.PythonPath, strings.Split(pp, ":")...)
	}
	return &lang, nil
}

func (l *Language) ParseFile(id string) (*language.FileInfo, error) {
	return python_grammar.Parse(id)
}
