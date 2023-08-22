package rust

import (
	"context"
	"fmt"
	"path"

	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/rust/rust_grammar"
	"github.com/gabotechs/dep-tree/internal/utils"
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

func CachedRustFile(ctx context.Context, id string) (context.Context, *rust_grammar.File, error) {
	cacheKey := language.FileCacheKey(id)
	if cached, ok := ctx.Value(cacheKey).(*rust_grammar.File); ok {
		return ctx, cached, nil
	} else {
		result, err := rust_grammar.Parse(id)
		if err != nil {
			return ctx, nil, err
		}
		ctx = context.WithValue(ctx, cacheKey, result)
		return ctx, result, err
	}
}

func MakeModTree(ctx context.Context, mainPath string, name string, parent *ModTree) (context.Context, *ModTree, error) {
	ctx, file, err := CachedRustFile(ctx, mainPath)
	if err != nil {
		return ctx, nil, err
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
					return ctx, nil, fmt.Errorf(`could not find mod "%s" in path "%s"`, stmt.Mod.Name, searchPath)
				}
				ctx, modTree.Children[stmt.Mod.Name], err = MakeModTree(ctx, modPath, stmt.Mod.Name, modTree)
				if err != nil {
					return ctx, nil, err
				}
			}
		}
	}

	return ctx, modTree, nil
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
