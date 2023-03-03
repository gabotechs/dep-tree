package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"dep-tree/internal/dep_tree"
	"dep-tree/internal/js"
	"dep-tree/internal/language"
	"dep-tree/internal/rust"
	"dep-tree/internal/tui"
	"dep-tree/internal/utils"
)

var jsonFormat bool

func run[T any, F any](
	ctx context.Context,
	entrypoint string,
	languageBuilder language.Builder[T, F],
) error {
	builder := language.ParserBuilder(languageBuilder)
	if jsonFormat {
		rendered, err := dep_tree.PrintStructured(ctx, entrypoint, builder)
		fmt.Println(rendered)
		return err
	}
	return tui.Loop(ctx, entrypoint, builder, nil)
}

func RenderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "render <path/to/entrypoint.ext>",
		Short: "Render the dependency tree starting from the provided entrypoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			entrypoint := args[0]

			switch {
			case utils.EndsWith(entrypoint, js.Extensions):
				return run(ctx, entrypoint, js.MakeJsLanguage)
			case utils.EndsWith(entrypoint, rust.Extensions):
				return run(ctx, entrypoint, rust.MakeRustLanguage)
			default:
				return fmt.Errorf("file \"%s\" not supported", entrypoint)
			}
		},
	}

	cmd.Flags().BoolVar(&jsonFormat, "json", false, "render the dependency try in a machine readable json format")

	return cmd
}
