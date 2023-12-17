package cmd

import (
	"context"
	"fmt"
	"github.com/gabotechs/dep-tree/internal/entropy"

	"github.com/spf13/cobra"

	"github.com/gabotechs/dep-tree/internal/js"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python"
	"github.com/gabotechs/dep-tree/internal/rust"
	"github.com/gabotechs/dep-tree/internal/utils"
)

func makeParserBuilder(ctx context.Context, entrypoint string) (context.Context, language.NodeParser, error) {
	switch {
	case utils.EndsWith(entrypoint, js.Extensions):
		return language.ParserBuilder(js.MakeJsLanguage, &cfg.Js)(ctx, entrypoint)
	case utils.EndsWith(entrypoint, rust.Extensions):
		return language.ParserBuilder(rust.MakeRustLanguage, &cfg.Rust)(ctx, entrypoint)
	case utils.EndsWith(entrypoint, python.Extensions):
		return language.ParserBuilder(python.MakePythonLanguage, &cfg.Python)(ctx, entrypoint)
	default:
		return ctx, nil, fmt.Errorf("file \"%s\" not supported", entrypoint)
	}
}

func EntropyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entropy",
		Short: "Renders a force-directed graph in the browser",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			entrypoint := args[0]

			err := loadConfig()
			if err != nil {
				return err
			}
			ctx, parser, err := makeParserBuilder(ctx, entrypoint)
			if err != nil {
				return err
			}
			ctx, err = entropy.Render(ctx, parser)
			return err
		},
	}

	return cmd
}
