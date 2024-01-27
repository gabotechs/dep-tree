package rust

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/gabotechs/dep-tree/internal/utils"
)

const cargoTomlFile = "Cargo.toml"

type localDependency struct {
	Path string
}

type CargoToml struct {
	// directory where the Cargo.toml file is located.
	path string
	// It's [dev-]dependencies.
	Dependencies map[string]localDependency
}

// readCargoToml parses a Cargo.toml file given its path or to the folder where it's placed.
var readCargoToml = utils.Cached1In2Out(func(path string) (*CargoToml, error) {
	path, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}
	fullPath := ""
	dir := ""
	if filepath.Base(path) != cargoTomlFile {
		fullPath = filepath.Join(path, cargoTomlFile)
		dir = path
	} else {
		fullPath = path
		dir = filepath.Dir(path)
	}
	var decoded struct {
		Dependencies    map[string]any `toml:"dependencies"`
		DevDependencies map[string]any `toml:"dev-dependencies"`
	}
	_, err = toml.DecodeFile(fullPath, &decoded)
	if err != nil {
		return nil, err
	}
	result := CargoToml{
		path:         dir,
		Dependencies: map[string]localDependency{},
	}
	for _, deps := range []map[string]any{decoded.DevDependencies, decoded.Dependencies} {
		for k, v := range deps {
			switch t := v.(type) {
			case map[string]any:
				switch tt := t["path"].(type) {
				case string:
					result.Dependencies[k] = localDependency{tt}
				}
			}
		}
	}
	return &result, nil
})

// findClosestCargoToml starts from a search path and goes up dir by dir
// until a Cargo.toml file is found. If one is found, it returns the
// parsed Cargo.toml file, if none is found, returns nil.
func findClosestCargoToml(searchPath string) (*CargoToml, error) {
	cargoTomlPath := filepath.Join(searchPath, "Cargo.toml")
	if utils.FileExists(cargoTomlPath) {
		return readCargoToml(cargoTomlPath)
	}
	nextSearchPath := filepath.Dir(searchPath)
	if nextSearchPath != searchPath {
		return findClosestCargoToml(nextSearchPath)
	} else {
		return nil, nil
	}
}

var searchPaths = []string{
	filepath.Join("src", "lib.rs"),
	filepath.Join("src", "main.rs"),
	"lib.rs",
	"main.rs",
}

// ModTree lazily builds the ModTree for this specific CargoToml.
//
// First call to this function needs to parse the ModTree, subsequent calls are cached.
func (c *CargoToml) ModTree() (*ModTree, error) {
	mainFile, err := c.MainFile()
	if err != nil {
		return nil, err
	}
	return MakeModTree(mainFile)
}

// MainFile retrieves the main file of the workspace (e.g. src/lib.rs or src/main.rs).
func (c *CargoToml) MainFile() (string, error) {
	for _, searchPath := range searchPaths {
		if p := filepath.Join(c.path, searchPath); utils.FileExists(p) {
			return p, nil
		}
	}
	return "", fmt.Errorf("main executable/library Rust file not found for cargo workspace in %s", c.path)
}
