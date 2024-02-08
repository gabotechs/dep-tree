package cmd

import (
	"fmt"

	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/spf13/cobra"

	"github.com/gabotechs/dep-tree/internal/check"
	"github.com/gabotechs/dep-tree/internal/config"
)

func CheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "check",
		Short:   "Checks that the dependency rules defined in the configuration file are not broken",
		GroupID: checkGroupId,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				configPath = config.DefaultConfigPath
			}
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			if len(cfg.Check.Entrypoints) == 0 {
				return fmt.Errorf(`config file "%s" has no entrypoints`, configPath)
			}
			parserBuilder, err := makeParserBuilder(cfg.Check.Entrypoints, cfg)
			if err != nil {
				return err
			}
			parser, err := parserBuilder(args)
			if err != nil {
				return err
			}
			return check.Check(parser, &cfg.Check, graph.NewStdErrCallbacks[*language.FileInfo]())
		},
	}
}
