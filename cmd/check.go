package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gabotechs/dep-tree/internal/config"
	"github.com/gabotechs/dep-tree/internal/js"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python"
	"github.com/gabotechs/dep-tree/internal/rust"
	"github.com/gabotechs/dep-tree/internal/utils"
)

func CheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Checks the that the current project matches the dependency rules defined in the configuration",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg, err := config.ParseConfig(configPath)
			if err != nil {
				return err
			}
			if len(cfg.Entrypoints) == 0 {
				return fmt.Errorf(`config file "%s" has no entrypoints`, configPath)
			}
			switch {
			case utils.EndsWith(cfg.Entrypoints[0], js.Extensions):
				return config.Check(ctx, language.ParserBuilder(js.MakeJsLanguage, &cfg.Js), cfg)
			case utils.EndsWith(cfg.Entrypoints[0], rust.Extensions):
				return config.Check(ctx, language.ParserBuilder(rust.MakeRustLanguage, &cfg.Rust), cfg)
			case utils.EndsWith(cfg.Entrypoints[0], python.Extensions):
				return config.Check(ctx, language.ParserBuilder(python.MakePythonLanguage, &cfg.Python), cfg)
			default:
				return fmt.Errorf("file \"%s\" not supported", cfg.Entrypoints[0])
			}
		},
	}

	return cmd
}
