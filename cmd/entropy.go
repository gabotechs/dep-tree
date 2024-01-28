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
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			files, err := filesFromArgs(args)
			if err != nil {
				return err
			}
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			parserBuilder, err := makeParserBuilder(files, cfg)
			if err != nil {
				return err
			}
			parser, err := parserBuilder(files)
			if err != nil {
				return err
			}
			err = entropy.Render(parser, files, entropy.RenderConfig{
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
