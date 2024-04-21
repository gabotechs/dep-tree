package dart

import (
	"os"
	"path/filepath"
	"sync"
)

// rootDir stores the found root directory to avoid repeated filesystem checks.
var rootDir string
var lock sync.Once

// findClosestDartRootDir finds the closest directory from the given path that contains a Dart project root indicator file.
// It caches the result after the first filesystem scan and reuses it for subsequent calls.
func findClosestDartRootDir(path string) string {
	lock.Do(func() {
		setRootDir(path)
	})
	return rootDir
}

// setRootDir performs the filesystem traversal to locate the root directory.
func setRootDir(path string) {
	var rootIndicatorFiles = []string{"pubspec.yaml", "pubspec.yml"}
	currentPath := path
	for {
		for _, file := range rootIndicatorFiles {
			if _, err := os.Stat(filepath.Join(currentPath, file)); err == nil {
				rootDir = currentPath
				return
			}
		}
		parentDir := filepath.Dir(currentPath)
		if parentDir == currentPath {
			panic("no Dart project root found. Make sure there is a pubspec.yaml or pubspec.yml in the project root.")
		}
		currentPath = parentDir
	}
}
