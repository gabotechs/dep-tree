package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"dep-tree/internal/dep_tree"
	"dep-tree/internal/js"
	"dep-tree/internal/language"
	"dep-tree/internal/tui"
)

func printStructured[T any](
	ctx context.Context,
	entrypoint string,
	parserBuilder func(string) (dep_tree.NodeParser[T], error),
) error {
	parser, err := parserBuilder(entrypoint)
	if err != nil {
		return err
	}
	_, dt, err := dep_tree.NewDepTree(ctx, parser)
	if err != nil {
		return err
	}
	output, err := dt.RenderStructured(parser.Display)
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

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
					return printStructured(ctx, entrypoint, language.ParserBuilder(js.MakeJsLanguage))
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
