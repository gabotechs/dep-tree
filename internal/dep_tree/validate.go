package dep_tree

import (
	"errors"
	"strings"

	"dep-tree/internal/config"
)

func (dt *DepTree[T]) Validate(cfg *config.Config) error {
	failures, err := cfg.Validate(dt.RootId, func(from string) []string {
		children := dt.Graph.Children(from)
		result := make([]string, len(children))
		for i, c := range children {
			result[i] = c.Id
		}
		return result
	})
	if err != nil {
		return err
	} else if len(failures) > 0 {
		return errors.New("the following dependencies are not allowed:\n" + strings.Join(failures, "\n"))
	}
	return nil
}
