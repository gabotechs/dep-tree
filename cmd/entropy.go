package cmd

import (
	"github.com/spf13/cobra"

	"github.com/gabotechs/dep-tree/internal/entropy"
)

var noBrowserOpen bool
var enableGui bool

func EntropyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "entropy",
		Short:   "(default) Renders a 3d force-directed graph in the browser",
		GroupID: renderGroupId,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			entrypoint := args[0]

			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			parserBuilder, err := makeParserBuilder(entrypoint, cfg)
			if err != nil {
				return err
			}
			ctx, parser, err := parserBuilder(ctx, entrypoint)
			if err != nil {
				return err
			}
			ctx, err = entropy.Render(ctx, parser, entropy.RenderConfig{
				NoOpen:    noBrowserOpen,
				EnableGui: enableGui,
			})
			return err
		},
	}

	cmd.Flags().BoolVar(&noBrowserOpen, "no-browser-open", false, "Disable the automatic browser open while rendering entropy")
	cmd.Flags().BoolVar(&enableGui, "enable-gui", false, "Enables a GUI for changing rendering settings")

	return cmd
}
