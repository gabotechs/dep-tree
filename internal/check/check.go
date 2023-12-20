package check

import (
	"context"
	"errors"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/utils"
)

func checkEntrypoint[T any](
	ctx context.Context,
	parserBuilder dep_tree.NodeParserBuilder[T],
	cfg *Config,
	entrypoint string,
) (context.Context, error) {
	ctx, parser, err := parserBuilder(ctx, entrypoint)
	if err != nil {
		return ctx, err
	}
	dt := dep_tree.NewDepTree(parser)
	root, err := dt.Root()
	if err != nil {
		return ctx, err
	}
	ctx, err = dt.LoadGraph(ctx)
	if err != nil {
		return ctx, err
	}
	failures, err := cfg.Validate(
		root.Id,
		func(from string) []string {
			toNodes := dt.Graph.FromId(from)
			result := make([]string, len(toNodes))
			for i, c := range toNodes {
				result[i] = c.Id
			}
			return result
		},
	)
	dt.LoadCycles()
	if !cfg.AllowCircularDependencies && dt.Cycles.Len() > 0 {
		for _, cycleId := range dt.Cycles.Keys() {
			cycle, _ := dt.Cycles.Get(cycleId)
			formattedCycleStack := make([]string, len(cycle.Stack))
			for i, el := range cycle.Stack {
				node := dt.Graph.Get(el)
				if node == nil {
					formattedCycleStack[i] = el
				} else {
					formattedCycleStack[i] = parser.Display(node)
				}
			}

			msg := "detected circular dependency: " + strings.Join(formattedCycleStack, " -> ")
			failures = append(failures, msg)
		}
	}
	if err != nil {
		return ctx, err
	} else if len(failures) > 0 {
		return ctx, errors.New("Check failed for entrypoint \"" + entrypoint + "\" the following dependencies are not allowed:\n" + strings.Join(failures, "\n"))
	}
	return ctx, nil
}

type Error []error

func (e Error) Error() string {
	msg := ""
	for _, err := range e {
		msg += err.Error()
		msg += "\n"
	}
	return msg
}

func Check[T any](
	ctx context.Context,
	parserBuilder dep_tree.NodeParserBuilder[T],
	cfg *Config,
) error {
	errorFlag := false
	errs := make([]error, len(cfg.Entrypoints))
	for i, entrypoint := range cfg.Entrypoints {
		ctx, errs[i] = checkEntrypoint(ctx, parserBuilder, cfg, path.Join(cfg.Path, entrypoint))
		if errs[i] != nil {
			errorFlag = true
		}
	}
	if errorFlag {
		return Error(errs)
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

func (c *Config) validate(
	start string,
	destinations func(from string) []string,
	seen map[string]bool,
) ([]string, error) {
	collectedErrors := make([]string, 0)

	if _, ok := seen[start]; ok {
		return collectedErrors, nil
	} else {
		seen[start] = true
	}

	for _, dest := range destinations(start) {
		from, to := c.rel(start), c.rel(dest)
		pass, err := c.Check(from, to)
		if err != nil {
			return nil, err
		} else if !pass {
			collectedErrors = append(collectedErrors, from+" -> "+to)
		}
		moreErrors, err := c.validate(dest, destinations, seen)
		if err != nil {
			return nil, err
		}
		collectedErrors = append(collectedErrors, moreErrors...)
	}
	return collectedErrors, nil
}

func (c *Config) Validate(start string, destinations func(from string) []string) ([]string, error) {
	return c.validate(start, destinations, map[string]bool{})
}
