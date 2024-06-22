package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/gabotechs/dep-tree/internal/config"
)

func ConfigCmd(_ func() (*config.Config, error)) *cobra.Command {
	return &cobra.Command{
		Use:     "config",
		Short:   "Generates a sample config in case that there's not already one present",
		Args:    cobra.ExactArgs(0),
		Aliases: []string{"init"},
		RunE: func(cmd *cobra.Command, args []string) error {
			path := config.DefaultConfigPath
			if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
				return os.WriteFile(path, []byte(config.SampleConfig), 0o600)
			} else {
				return fmt.Errorf("cannot generate config file, as one already exists in %s", path)
			}
		},
	}
}
