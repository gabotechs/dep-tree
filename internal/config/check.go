package config

import (
	"context"
	"errors"
	"path"
	"strings"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
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
	ctx, dt, err := dep_tree.NewDepTree(ctx, parser)
	if err != nil {
		return ctx, err
	}
	failures, err := cfg.Validate(
		dt.RootId,
		func(from string) []string {
			children := dt.Graph.Children(from)
			result := make([]string, len(children))
			for i, c := range children {
				result[i] = c.Id
			}
			return result
		},
	)
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

type CheckError []error

func (e CheckError) Error() string {
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
		return CheckError(errs)
	}
	return nil
}
