package check

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/utils"
)

func Check[T any](
	parser graph.NodeParser[T],
	display func(node *graph.Node[T]) string,
	cfg *Config,
	callbacks graph.LoadCallbacks[T],
) error {
	// 1. build the graph.
	files := make([]string, len(cfg.Entrypoints))
	for i, file := range cfg.Entrypoints {
		files[i] = filepath.Join(cfg.Path, file)
	}

	g := graph.NewGraph[T]()
	err := g.Load(files, parser, callbacks)
	if err != nil {
		return err
	}
	// 2. Check for rule violations in the graph.
	failures := make([]string, 0)
	for _, node := range g.AllNodes() {
		for _, dep := range g.FromId(node.Id) {
			from, to := cfg.rel(node.Id), cfg.rel(dep.Id)
			pass, err := cfg.Check(from, to)
			if err != nil {
				return err
			} else if !pass {
				failures = append(failures, from+" -> "+to)
			}
		}
	}
	// 3. Check for cycles.
	cycles := g.RemoveElementaryCycles()
	if !cfg.AllowCircularDependencies {
		for _, cycle := range cycles {
			formattedCycleStack := make([]string, len(cycle.Stack))
			for i, el := range cycle.Stack {
				if node := g.Get(el); node != nil {
					formattedCycleStack[i] = display(node)
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
