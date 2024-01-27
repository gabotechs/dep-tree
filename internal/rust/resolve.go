package rust

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// filePathToModChain builds the mod chain for the file located in the filePath argument
// based on its relative position to the main file (src/lib.rs, src/main.rs...)
func (l *Language) filePathToModChain(filePath string, mainFile string) ([]string, error) {
	mainFile, err := filepath.Abs(mainFile)
	if err != nil {
		return nil, err
	}
	filePath, err = filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}
	if filePath == mainFile {
		return []string{}, nil
	}
	root := filepath.Dir(mainFile)
	rel, err := filepath.Rel(root, filePath)
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

// resolve resolves which file is imported in a `use` statement.
//
// it receives the list of path slices in the `use` statement (e.g. `use foo::bar::baz` -> ["foo", "bar", "baz"])
// and the path of the file that contains the `use` statement, and returns the absolute imported path.
//
// we need the last `filePath` argument because:
// - we need to check what is the closest Cargo.toml file
// - we need
func (l *Language) resolve(pathSlices []string, filePath string) (string, error) {
	if len(pathSlices) == 0 {
		return filePath, nil
	}

	first := pathSlices[0]
	var modSearch []string

	cargoToml, err := findClosestCargoToml(filePath)
	if err != nil {
		return "", err
	}
	mainFile, err := cargoToml.MainFile()
	if err != nil {
		return "", err
	}
	modTree, err := cargoToml.ModTree()
	if err != nil {
		return "", err
	}

	if first == crate {
		modSearch = pathSlices[1:]
	} else {
		mods, err := l.filePathToModChain(filePath, mainFile)
		if err != nil {
			return "", err
		}
		mods = append(mods, pathSlices...)
		modSearch = mods
	}

	mod := modTree.Search(modSearch)
	switch {
	case mod == nil && (first == self || first == super || first == crate):
		return "", fmt.Errorf("could not find mod chain %s in the projects mod tree", strings.Join(modSearch, " -> "))
	case mod == nil:
		return "", nil
	default:
		return mod.Path, nil
	}
}
