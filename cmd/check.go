package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"dep-tree/internal/config"
	"dep-tree/internal/js"
	"dep-tree/internal/language"
	"dep-tree/internal/python"
	"dep-tree/internal/rust"
	"dep-tree/internal/utils"
)

var configPath string

func CheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Checks the that the current project matches the dependency rules defined in the configuration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg, err := config.ParseConfig(configPath)
			if err != nil {
				return fmt.Errorf("could not parse config file %s: %w", configPath, err)
			}
			switch {
			case utils.EndsWith(cfg.Entrypoints[0], js.Extensions):
				return config.Check(ctx, language.ParserBuilder(js.MakeJsLanguage), cfg)
			case utils.EndsWith(cfg.Entrypoints[0], rust.Extensions):
				return config.Check(ctx, language.ParserBuilder(rust.MakeRustLanguage), cfg)
			case utils.EndsWith(cfg.Entrypoints[0], python.Extensions):
				return config.Check(ctx, language.ParserBuilder(python.MakePythonLanguage), cfg)
			default:
				return fmt.Errorf("file \"%s\" not supported", cfg.Entrypoints[0])
			}
		},
	}

	cmd.Flags().StringVar(&configPath, "config", ".dep-tree.yml", "path to dep-tree's config file")

	return cmd
}
