package dep_tree

import (
	"context"
	"errors"
	"strings"

	"dep-tree/internal/config"
)

func (dt *DepTree[T]) Validate(
	ctx context.Context,
	cfg *config.Config,
) (context.Context, error) {
	failures, err := cfg.Validate(dt.RootId, func(from string) []string {
		children := dt.Graph.Children(from)
		result := make([]string, len(children))
		for i, c := range children {
			result[i] = c.Id
		}
		return result
	})
	if err != nil {
		return ctx, err
	} else if len(failures) > 0 {
		return ctx, errors.New("the following dependencies are not allowed:\n" + strings.Join(failures, "\n"))
	}
	return ctx, nil
}
