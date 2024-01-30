package check

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/utils"
)

func Check[T any](parser dep_tree.NodeParser[T], cfg *Config) error {
	dt := dep_tree.NewDepTree(parser, cfg.Entrypoints).WithStdErrLoader()
	err := dt.LoadGraph()
	if err != nil {
		return err
	}
	// 1. Check for rule violations in the graph.
	failures := make([]string, 0)
	for _, node := range dt.Graph.AllNodes() {
		for _, dep := range dt.Graph.FromId(node.Id) {
			from, to := cfg.rel(node.Id), cfg.rel(dep.Id)
			pass, err := cfg.Check(from, to)
			if err != nil {
				return err
			} else if !pass {
				failures = append(failures, from+" -> "+to)
			}
		}
	}
	// 2. Check for cycles.
	dt.LoadCycles()
	if !cfg.AllowCircularDependencies {
		for el := dt.Cycles.Front(); el != nil; el = el.Next() {
			formattedCycleStack := make([]string, len(el.Value.Stack))
			for i, el := range el.Value.Stack {
				if node := dt.Graph.Get(el); node != nil {
					formattedCycleStack[i] = parser.Display(node)
				} else {
					formattedCycleStack[i] = el
				}
			}

			msg := "detected circular dependency: " + strings.Join(formattedCycleStack, " -> ")
			failures = append(failures, msg)
		}
	}
	if len(failures) > 0 {
		return errors.New("Check failed, the following dependencies are not allowed:\n" + strings.Join(failures, "\n"))
	}
	return nil
}

func (c *Config) whiteListCheck(from, to string) (bool, error) {
	for k, v := range c.WhiteList {
		doesMatch, err := utils.GlobstarMatch(k, from)
		if err != nil {
			return false, err
		}
		if doesMatch {
			for _, dest := range v {
				shouldPass, err := utils.GlobstarMatch(dest, to)
				if err != nil {
					return false, err
				}
				if shouldPass {
					return true, nil
				}
			}
			return false, nil
		}
	}
	return true, nil
}

func (c *Config) blackListCheck(from, to string) (bool, error) {
	for k, v := range c.BlackList {
		doesMatch, err := utils.GlobstarMatch(k, from)
		if err != nil {
			return false, err
		}
		if doesMatch {
			for _, dest := range v {
				shouldReject, err := utils.GlobstarMatch(dest, to)
				if err != nil {
					return false, err
				}
				if shouldReject {
					return false, nil
				}
			}
		}
	}

	return true, nil
}

func (c *Config) Check(from, to string) (bool, error) {
	pass, err := c.blackListCheck(from, to)
	if err != nil || !pass {
		return pass, err
	}
	return c.whiteListCheck(from, to)
}

func (c *Config) rel(p string) string {
	relPath, err := filepath.Rel(c.Path, p)
	if err != nil {
		return p
	}
	return relPath
}
