package cmd

import (
	"context"
	"fmt"
	"github.com/gabotechs/dep-tree/internal/config"
	"os"

	"github.com/spf13/cobra"

	"github.com/gabotechs/dep-tree/internal/dep_tree"
	"github.com/gabotechs/dep-tree/internal/js"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python"
	"github.com/gabotechs/dep-tree/internal/rust"
	"github.com/gabotechs/dep-tree/internal/tui"
	"github.com/gabotechs/dep-tree/internal/utils"
)

var jsonFormat bool

func run[T any, F any, C any](
	ctx context.Context,
	entrypoint string,
	languageBuilder language.Builder[T, F, C],
	cfg C,
) error {
	builder := language.ParserBuilder(languageBuilder, cfg)
	if jsonFormat {
		rendered, err := dep_tree.PrintStructured(ctx, entrypoint, builder)
		fmt.Println(rendered)
		return err
	}
	return tui.Loop(ctx, entrypoint, builder, nil, true, nil)
}

func runRender(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	entrypoint := args[0]

	cfg, err := config.ParseConfig(configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	switch {
	case utils.EndsWith(entrypoint, js.Extensions):
		return run(ctx, entrypoint, js.MakeJsLanguage, &cfg.Js)
	case utils.EndsWith(entrypoint, rust.Extensions):
		return run(ctx, entrypoint, rust.MakeRustLanguage, &cfg.Rust)
	case utils.EndsWith(entrypoint, python.Extensions):
		return run(ctx, entrypoint, python.MakePythonLanguage, &cfg.Python)
	default:
		return fmt.Errorf("file \"%s\" not supported", entrypoint)
	}
}

func RenderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "render <path/to/entrypoint.ext>",
		Short: "[default] Render the dependency tree starting from the provided entrypoint",
		Args:  cobra.ExactArgs(1),
		RunE:  runRender,
	}

	cmd.Flags().BoolVar(&jsonFormat, "json", false, "render the dependency try in a machine readable json format")

	return cmd
}
