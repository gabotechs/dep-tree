package cmd

import (
	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/language"
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
			lang, err := inferLang(files, cfg)
			if err != nil {
				return err
			}
			parser := language.NewParser(lang)
			parser.UnwrapProxyExports = cfg.UnwrapExports
			parser.Exclude = cfg.Exclude
			err = entropy.Render(files, parser, entropy.RenderConfig{
				NoOpen:        noBrowserOpen,
				EnableGui:     enableGui,
				LoadCallbacks: graph.NewStdErrCallbacks[*language.FileInfo](),
			})
			return err
		},
	}

	cmd.Flags().BoolVar(&noBrowserOpen, "no-browser-open", false, "Disable the automatic browser open while rendering entropy")
	cmd.Flags().BoolVar(&enableGui, "enable-gui", false, "Enables a GUI for changing rendering settings")

	return cmd
}
