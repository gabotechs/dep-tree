package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"dep-tree/internal/js"
	"dep-tree/internal/tui"
)

func RenderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "render <path/to/entrypoint.ext>",
		Short: "Render the dependency tree starting from the provided entrypoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			entrypoint := args[0]

			if endsWith(entrypoint, js.Extensions) {
				return tui.Loop(
					ctx,
					entrypoint,
					js.MakeJsParser,
					nil,
				)
			} else {
				return fmt.Errorf("file \"%s\" not supported", entrypoint)
			}
		},
	}

	return cmd
}
