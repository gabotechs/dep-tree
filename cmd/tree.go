package cmd

import (
	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/tree"
	"github.com/spf13/cobra"

	"github.com/gabotechs/dep-tree/internal/tui"
)

var jsonFormat bool

func TreeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tree",
		Short:   "Render the dependency tree starting from the provided entrypoint",
		Args:    cobra.MinimumNArgs(1),
		GroupID: renderGroupId,
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

			if jsonFormat {
				t, err := tree.NewTree[*language.FileInfo](
					files,
					parser,
					relPathDisplay,
					graph.NewStdErrCallbacks[*language.FileInfo](relPathDisplay),
				)
				if err != nil {
					return err
				}

				rendered, err := t.RenderStructured()
				cmd.Println(rendered)
				return err
			} else {
				return tui.Loop[*language.FileInfo](
					files,
					parser,
					relPathDisplay,
					nil,
					true,
					nil,
					graph.NewStdErrCallbacks[*language.FileInfo](relPathDisplay))
			}
		},
	}

	cmd.Flags().BoolVar(&jsonFormat, "json", false, "render the dependency tree in a machine readable json format")

	return cmd
}
