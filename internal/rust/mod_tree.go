package rust

import (
	"fmt"
	"path"
	"sync"

	"dep-tree/internal/rust/rust_grammar"
	"dep-tree/internal/utils"
)

type ModTree struct {
	Name     string
	Path     string
	Parent   *ModTree
	Children map[string]*ModTree
}

const self = "self"
const crate = "crate"
const super = "super"

// TODO: propagating context here is very messy.
var fileCache = make(map[string]*rust_grammar.File)
var fileCacheMutex = sync.Mutex{}

func MakeModTree(mainPath string, name string, parent *ModTree) (*ModTree, error) {
	var file *rust_grammar.File
	var ok bool
	var err error
	if file, ok = fileCache[mainPath]; !ok {
		file, err = rust_grammar.Parse(mainPath)
		if err != nil {
			return nil, err
		}
		fileCacheMutex.Lock()
		fileCache[mainPath] = file
		fileCacheMutex.Unlock()
	}

	var searchPath string
	if path.Base(mainPath) == name+".rs" {
		searchPath = path.Join(path.Dir(mainPath), name)
	} else {
		searchPath = path.Dir(mainPath)
	}

	modTree := &ModTree{
		Name:     name,
		Path:     mainPath,
		Parent:   parent,
		Children: make(map[string]*ModTree),
	}

	for _, stmt := range file.Statements {
		if stmt.Mod != nil {
			if stmt.Mod.Local {
				modTree.Children[stmt.Mod.Name] = &ModTree{
					Name: stmt.Mod.Name,
					Path: mainPath,
				}
			} else {
				var modPath string
				if p := path.Join(searchPath, stmt.Mod.Name+".rs"); utils.FileExists(p) {
					modPath = p
				} else if p = path.Join(searchPath, stmt.Mod.Name, "mod.rs"); utils.FileExists(p) {
					modPath = p
				} else {
					return nil, fmt.Errorf(`coult not find mod "%s" in path "%s"`, stmt.Mod.Name, searchPath)
				}
				modTree.Children[stmt.Mod.Name], err = MakeModTree(modPath, stmt.Mod.Name, modTree)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return modTree, nil
}

func (m *ModTree) Search(modChain []string) *ModTree {
	current := m
	for _, mod := range modChain {
		if mod == self {
			continue
		} else if mod == super {
			if current.Parent == nil {
				return nil
			} else {
				current = current.Parent
			}
		} else if child, ok := current.Children[mod]; ok {
			current = child
		} else {
			return nil
		}
	}
	return current
}
