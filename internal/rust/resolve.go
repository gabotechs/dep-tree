package rust

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// filePathToModChain builds the mod chain for the file located in the filePath argument
// based on its relative position to the main file (src/lib.rs, src/main.rs...)
func filePathToModChain(filePath string, mainFile string) ([]string, error) {
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
func resolve(pathSlices []string, filePath string) (string, error) {
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

	var mod *ModTree

	if first == crate {
		// Referring to a mod in this workspace, and the path slices are relative to the root of the workspace.
		mod = modTree.Search(pathSlices[1:])
	} else if first == self || first == super {
		// Referring to a mod in this workspace, and the path slices are relative to the evaluated file.
		mods, err := filePathToModChain(filePath, mainFile)
		if err != nil {
			return "", err
		}
		mod = modTree.Search(append(mods, pathSlices...))
	} else if workspace, ok := cargoToml.Dependencies[first]; ok {
		// Referring to a mod in another workspace
		workspaceRoot, err := filepath.Abs(filepath.Join(cargoToml.path, workspace.Path))
		if err != nil {
			return "", fmt.Errorf("could not find workspace %s relative to %s: %w", workspace.Path, cargoToml.path, err)
		}
		cargoToml, err = readCargoToml(workspaceRoot)
		if err != nil {
			return "", err
		}
		modTree, err := cargoToml.ModTree()
		if err != nil {
			return "", fmt.Errorf("could not create mod tree for workspace %s: %w", workspace.Path, err)
		}
		mod = modTree.Search(pathSlices[1:])
	} else {
		return "", nil
	}

	if mod == nil {
		return "", fmt.Errorf("could not find mod chain %s in the projects mod tree", strings.Join(modSearch, " -> "))
	}
	return mod.Path, nil
}
