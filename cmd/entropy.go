package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gabotechs/dep-tree/internal/entropy"
	"github.com/gabotechs/dep-tree/internal/js"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python"
	"github.com/gabotechs/dep-tree/internal/rust"
	"github.com/gabotechs/dep-tree/internal/utils"
)

func runEntropy(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	entrypoint := args[0]

	err := loadConfig()
	if err != nil {
		return err
	}
	switch {
	case utils.EndsWith(entrypoint, js.Extensions):
		ctx, parser, err := language.ParserBuilder(js.MakeJsLanguage, &cfg.Js)(ctx, entrypoint)
		if err != nil {
			return err
		}
		_, err = entropy.Render(ctx, parser)
		return err
	case utils.EndsWith(entrypoint, rust.Extensions):
		ctx, parser, err := language.ParserBuilder(rust.MakeRustLanguage, &cfg.Rust)(ctx, entrypoint)
		if err != nil {
			return err
		}
		_, err = entropy.Render(ctx, parser)
		return err
	case utils.EndsWith(entrypoint, python.Extensions):
		ctx, parser, err := language.ParserBuilder(python.MakePythonLanguage, &cfg.Python)(ctx, entrypoint)
		if err != nil {
			return err
		}
		_, err = entropy.Render(ctx, parser)
		return err
	default:
		return fmt.Errorf("file \"%s\" not supported", entrypoint)
	}
}

func EntropyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entropy",
		Short: "Renders a force-directed graph in the browser",
		Args:  cobra.ExactArgs(1),
		RunE:  runEntropy,
	}

	return cmd
}
