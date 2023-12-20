package cmd

import (
	"github.com/spf13/cobra"

	"github.com/gabotechs/dep-tree/internal/config"

	"github.com/gabotechs/dep-tree/internal/entropy"
)

func EntropyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entropy",
		Short: "Renders a force-directed graph in the browser",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			entrypoint := args[0]

			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			if cfg == nil {
				cfg = &config.Config{}
			}
			if cfg.FollowReExports == nil {
				v := false
				cfg.FollowReExports = &v
			}
			parserBuilder, err := makeParserBuilder(entrypoint, cfg)
			if err != nil {
				return err
			}
			ctx, parser, err := parserBuilder(ctx, entrypoint)
			if err != nil {
				return err
			}
			ctx, err = entropy.Render(ctx, parser)
			return err
		},
	}

	return cmd
}
