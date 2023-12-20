package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/gabotechs/dep-tree/internal/config"
)

func ConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "config",
		Short:   "Generates a sample config in case that there's not already one present",
		Args:    cobra.ExactArgs(0),
		Aliases: []string{"init"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				configPath = config.DefaultConfigPath
			}
			if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
				return os.WriteFile(configPath, []byte(config.SampleConfig), os.ModePerm)
			} else {
				return errors.New("Cannot generate config file, as one already exists in " + configPath)
			}
		},
	}
}
