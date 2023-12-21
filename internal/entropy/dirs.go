package entropy

import (
	"path"
	"strings"
)

func dirs(dir string) []string {
	var result []string
	for strings.Contains(dir, "/") {
		if strings.HasSuffix(dir, "/") {
			dir = dir[:len(dir)-1]
		}
		if strings.HasPrefix(dir, "/") {
			dir = dir[1:]
		}
		result = append(result, dir)
		dir = path.Dir(dir)
	}
	if dir != "" && dir != "." {
		result = append(result, dir)
	}
	return result
}
