package cpp

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/language"
)

var Extensions = []string{
	"cpp", "cc", "cxx", "c++",
	"hpp", "hh", "hxx", "h++", "h",
}

type Language struct {
	Cfg *Config
}

func NewLanguage(cfg *Config) *Language {
	if cfg == nil {
		cfg = &Config{}
	}
	return &Language{Cfg: cfg}
}

func (l *Language) ParseFile(path string) (*language.FileInfo, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	file, err := ParseCppFile(string(content))
	if err != nil {
		return nil, err
	}

	currentDir, _ := os.Getwd()
	relPath, _ := filepath.Rel(currentDir, path)

	return &language.FileInfo{
		Content: file,
		Loc:     bytes.Count(content, []byte("\n")),
		Size:    len(content),
		AbsPath: path,
		RelPath: relPath,
	}, nil
}

func (l *Language) ParseImports(file *language.FileInfo) (*language.ImportsResult, error) {
	var result language.ImportsResult

	cppFile, ok := file.Content.(*File)
	if !ok {
		return &result, nil
	}

	for _, include := range cppFile.Includes {
		if include.IsSystem {
			continue
		}

		absPath := l.resolveIncludePath(file.AbsPath, include.Header)
		if absPath != "" {
			result.Imports = append(result.Imports, language.ImportEntry{
				All:     true,
				AbsPath: absPath,
			})
		}
	}

	return &result, nil
}

func (l *Language) ParseExports(file *language.FileInfo) (*language.ExportsResult, error) {
	var result language.ExportsResult

	// For C++, determining exports is complex as it depends on:
	// - Public class/struct members
	// - Free functions
	// - Global variables
	// - Template definitions
	// For now, we'll treat header files as exporting everything
	// and source files as exporting nothing by default

	if l.isHeaderFile(file.AbsPath) {
		// Header files export all their content
		result.Exports = append(result.Exports, language.ExportEntry{
			All:     true,
			AbsPath: file.AbsPath,
		})
	}

	return &result, nil
}

func (l *Language) resolveIncludePath(sourceFile, includePath string) string {
	// If the include path is absolute, return it as-is
	if filepath.IsAbs(includePath) {
		return includePath
	}

	// try relative to the source file directory
	sourceDir := filepath.Dir(sourceFile)
	resolvedPath := filepath.Join(sourceDir, includePath)

	if _, err := os.Stat(resolvedPath); err == nil {
		abs, _ := filepath.Abs(resolvedPath)
		return abs
	}

	if !hasExtension(includePath) {
		for _, ext := range []string{".h", ".hpp", ".hxx", ".h++"} {
			testPath := resolvedPath + ext
			if _, err := os.Stat(testPath); err == nil {
				abs, _ := filepath.Abs(testPath)
				return abs
			}
		}
	}

	// try relative to project root
	projectRoots := l.findProjectRoots(sourceDir)
	for _, root := range projectRoots {
		testPath := filepath.Join(root, includePath)
		if _, err := os.Stat(testPath); err == nil {
			abs, _ := filepath.Abs(testPath)
			return abs
		}

		// Try with extensions
		if !hasExtension(includePath) {
			for _, ext := range []string{".h", ".hpp", ".hxx", ".h++"} {
				testPathWithExt := testPath + ext
				if _, err := os.Stat(testPathWithExt); err == nil {
					abs, _ := filepath.Abs(testPathWithExt)
					return abs
				}
			}
		}
	}

	// try supporting configured include directories
	for _, includeDir := range l.Cfg.IncludeDirs {
		var testPath string
		if filepath.IsAbs(includeDir) {
			testPath = filepath.Join(includeDir, includePath)
		} else {
			testPath = filepath.Join(sourceDir, includeDir, includePath)
		}

		if _, err := os.Stat(testPath); err == nil {
			abs, _ := filepath.Abs(testPath)
			return abs
		}

		// Try with extensions
		if !hasExtension(includePath) {
			for _, ext := range []string{".h", ".hpp", ".hxx", ".h++"} {
				testPathWithExt := testPath + ext
				if _, err := os.Stat(testPathWithExt); err == nil {
					abs, _ := filepath.Abs(testPathWithExt)
					return abs
				}
			}
		}
	}

	return ""
}

func (l *Language) findProjectRoots(startDir string) []string {
	var roots []string
	currentDir := startDir

	// Common project root indicators
	rootIndicators := []string{
		"CMakeLists.txt", "Makefile", "SConstruct", // Build files
		".git", ".hg", ".svn", // Version control
		"package.json", "Cargo.toml", "go.mod", // Language-specific
		"README.md", "README.txt", // Documentation
	}

	for {
		for _, indicator := range rootIndicators {
			if _, err := os.Stat(filepath.Join(currentDir, indicator)); err == nil {
				roots = append(roots, currentDir)
				break
			}
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir

		if len(roots) >= 3 {
			break
		}
	}

	return roots
}

func (l *Language) isHeaderFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	headerExts := []string{".h", ".hpp", ".hh", ".hxx", ".h++"}
	for _, headerExt := range headerExts {
		if ext == headerExt {
			return true
		}
	}
	return false
}

func hasExtension(path string) bool {
	return filepath.Ext(path) != ""
}
