package cmd

import (
	"fmt"

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

			parserBuilder, err := makeParserBuilder(files, cfg)
			if err != nil {
				return err
			}

			if jsonFormat {
				parser, err := parserBuilder(files)
				if err != nil {
					return err
				}

				depTree, err := tree.NewTree(files, parser, graph.NewStdErrCallbacks[*language.FileInfo]())
				if err != nil {
					return err
				}

				rendered, err := depTree.RenderStructured()
				fmt.Println(rendered)
				return err
			} else {
				return tui.Loop(files, parserBuilder, nil, true, nil, graph.NewStdErrCallbacks[*language.FileInfo]())
			}
		},
	}

	cmd.Flags().BoolVar(&jsonFormat, "json", false, "render the dependency tree in a machine readable json format")

	return cmd
}
