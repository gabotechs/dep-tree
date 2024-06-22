package cmd

import (
	"github.com/spf13/cobra"

	"github.com/gabotechs/dep-tree/internal/config"
	"github.com/gabotechs/dep-tree/internal/entropy"
	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/language"
)

func EntropyCmd(cfgF func() (*config.Config, error)) *cobra.Command {
	var noBrowserOpen bool
	var enableGui bool
	var renderPath string

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
			cfg, err := cfgF()
			if err != nil {
				return err
			}
			lang, err := inferLang(files, cfg)
			if err != nil {
				return err
			}
			parser := language.NewParser(lang)
			applyConfigToParser(parser, cfg)

			err = entropy.Render(files, parser, entropy.RenderConfig{
				NoOpen:        noBrowserOpen,
				EnableGui:     enableGui,
				LoadCallbacks: graph.NewStdErrCallbacks[*language.FileInfo](relPathDisplay),
				RenderPath:    renderPath,
			})
			return err
		},
	}

	cmd.Flags().BoolVar(&noBrowserOpen, "no-browser-open", false, "Disable the automatic browser open while rendering entropy")
	cmd.Flags().BoolVar(&enableGui, "enable-gui", false, "Enables a GUI for changing rendering settings")
	cmd.Flags().StringVar(&renderPath, "render-path", "", "Sets the output path of the rendered html file")

	return cmd
}
