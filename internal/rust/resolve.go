package rust

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (l *Language) filePathToModChain(dir string) ([]string, error) {
	if dir == l.ProjectEntrypoint {
		return []string{}, nil
	}
	root := filepath.Dir(l.ProjectEntrypoint)
	rel, err := filepath.Rel(root, dir)
	if err != nil {
		return nil, err
	}
	slices := strings.Split(rel, string(os.PathSeparator))
	filteredSlices := make([]string, 0)
	for _, slice := range slices {
		switch {
		case slice == ".":
			continue
		case strings.HasSuffix(slice, ".rs"):
			slice = slice[:len(slice)-3]
			if slice != "mod" {
				filteredSlices = append(filteredSlices, slice)
			}
		default:
			filteredSlices = append(filteredSlices, slice)
		}
	}
	return filteredSlices, nil
}

func (l *Language) resolve(pathSlices []string, filePath string) (string, error) {
	if len(pathSlices) == 0 {
		return filePath, nil
	}

	first := pathSlices[0]
	var modSearch []string

	if first == crate {
		modSearch = pathSlices[1:]
	} else {
		mods, err := l.filePathToModChain(filePath)
		if err != nil {
			return "", err
		}
		mods = append(mods, pathSlices...)
		modSearch = mods
	}

	mod := l.ModTree.Search(modSearch)
	switch {
	case mod == nil && (first == self || first == super || first == crate):
		return "", fmt.Errorf("could not find mod chain %s in the projects mod tree", strings.Join(modSearch, " -> "))
	case mod == nil:
		return "", nil
	default:
		return mod.Path, nil
	}
}
