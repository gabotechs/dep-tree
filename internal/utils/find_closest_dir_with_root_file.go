package utils

import "path/filepath"

type SourcesRoot struct {
	FoundFile string
	AbsDir    string
}

func _findClosestDirWithRootFile(searchPath string, rootFiles []string) *SourcesRoot {
	for _, rootFile := range rootFiles {
		if FileExists(filepath.Join(searchPath, rootFile)) {
			return &SourcesRoot{
				FoundFile: rootFile,
				AbsDir:    searchPath,
			}
		}
	}
	nextSearchPath := filepath.Dir(searchPath)

	if nextSearchPath != searchPath {
		return _findClosestDirWithRootFile(nextSearchPath, rootFiles)
	} else {
		return nil
	}
}

func MakeCachedFindClosestDirWithRootFile(rootFiles []string) func(string) *SourcesRoot {
	f := func(searchPath string) *SourcesRoot {
		return _findClosestDirWithRootFile(searchPath, rootFiles)
	}
	return Cached1In1Out(f)
}
