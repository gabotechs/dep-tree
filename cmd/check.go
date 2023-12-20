package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gabotechs/dep-tree/internal/check"
	"github.com/gabotechs/dep-tree/internal/config"
)

func CheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Checks the that the current project matches the dependency rules defined in the configuration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if configPath == "" {
				configPath = config.DefaultConfigPath
			}
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			if len(cfg.Check.Entrypoints) == 0 {
				return fmt.Errorf(`config file "%s" has no entrypoints`, configPath)
			}
			parserBuilder, err := makeParserBuilder(cfg.Check.Entrypoints[0], cfg)
			if err != nil {
				return err
			}
			return check.Check(ctx, parserBuilder, &cfg.Check)
		},
	}
}
