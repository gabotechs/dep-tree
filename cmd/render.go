package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"dep-tree/internal/dep_tree"
	"dep-tree/internal/js"
	"dep-tree/internal/language"
	"dep-tree/internal/tui"
)

func RenderCmd() *cobra.Command {
	var jsonFormat bool

	cmd := &cobra.Command{
		Use:   "render <path/to/entrypoint.ext>",
		Short: "Render the dependency tree starting from the provided entrypoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			entrypoint := args[0]

			if endsWith(entrypoint, js.Extensions) {
				if jsonFormat {
					rendered, err := dep_tree.PrintStructured(
						ctx,
						entrypoint,
						language.ParserBuilder(js.MakeJsLanguage),
					)
					fmt.Println(rendered)
					return err
				}

				return tui.Loop(
					ctx,
					entrypoint,
					language.ParserBuilder(js.MakeJsLanguage),
					nil,
				)
			} else {
				return fmt.Errorf("file \"%s\" not supported", entrypoint)
			}
		},
	}

	cmd.Flags().BoolVar(&jsonFormat, "json", false, "render the dependency try in a machine readable json format")

	return cmd
}
