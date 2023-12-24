package utils

import "os"

func _FileExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

var FileExists = Cached(_FileExists)

func _DirExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

var DirExists = Cached(_DirExists)
