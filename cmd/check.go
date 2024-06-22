package cmd

import (
	"errors"
	"fmt"

	"github.com/gabotechs/dep-tree/internal/config"
	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/spf13/cobra"

	"github.com/gabotechs/dep-tree/internal/check"
)

func CheckCmd(cfgF func() (*config.Config, error)) *cobra.Command {
	return &cobra.Command{
		Use:     "check",
		Short:   "Checks that the dependency rules defined in the configuration file are not broken",
		GroupID: checkGroupId,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := cfgF()
			if err != nil {
				return err
			}
			if cfg.Source == "default" {
				return errors.New("when using the `check` subcommand, a .dep-tree.yml file must be provided, you can create one sample .dep-tree.yml file executing `dep-tree config` in your terminal")
			}

			if len(cfg.Check.Entrypoints) == 0 {
				return fmt.Errorf(`config file "%s" has no entrypoints`, cfg.Path)
			}
			lang, err := inferLang(cfg.Check.Entrypoints, cfg)
			if err != nil {
				return err
			}
			parser := language.NewParser(lang)
			applyConfigToParser(parser, cfg)

			return check.Check[*language.FileInfo](
				parser,
				relPathDisplay,
				&cfg.Check,
				graph.NewStdErrCallbacks[*language.FileInfo](relPathDisplay),
			)
		},
	}
}
