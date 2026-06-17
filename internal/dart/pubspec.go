package dart

import (
	"os"
	"path/filepath"

	"github.com/gabotechs/dep-tree/internal/utils"
	"gopkg.in/yaml.v3"
)

const pubspecFile = "pubspec.yaml"

// pubspec holds the relevant information parsed from a Dart `pubspec.yaml` file.
type pubspec struct {
	// Name is the value of the `name` attribute, which is the name of the Dart/Flutter package.
	Name string
	// absDir is the absolute path to the directory where the pubspec.yaml file is located.
	absDir string
}

var readPubspec = utils.Cached1In1OutErr(func(path string) (*pubspec, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Name string `yaml:"name"`
	}
	if err := yaml.Unmarshal(content, &parsed); err != nil {
		return nil, err
	}
	return &pubspec{
		Name:   parsed.Name,
		absDir: filepath.Dir(path),
	}, nil
})

var findClosestDirWithPubspec = utils.MakeCachedFindClosestDirWithRootFile([]string{pubspecFile})

// findClosestPubspec walks up the directory tree starting from searchPath until it
// finds a directory containing a pubspec.yaml file, and returns the parsed pubspec.
// If none is found, it returns nil.
func findClosestPubspec(searchPath string) *pubspec {
	root := findClosestDirWithPubspec(searchPath)
	if root == nil {
		return nil
	}
	result, err := readPubspec(filepath.Join(root.AbsDir, root.FoundFile))
	if err != nil {
		return nil
	}
	return result
}
