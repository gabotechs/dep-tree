package cmd

import (
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/spf13/cobra"
)

func ExplainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "explain",
		Short:   "Shows all the dependencies between two parts of the code",
		GroupID: explainGroupId,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromFiles, err := filesFromArgs([]string{args[0]})
			if err != nil {
				return err
			}
			toFiles, err := filesFromArgs([]string{args[0]})
			if err != nil {
				return err
			}

			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			lang, err := inferLang(fromFiles, cfg)
			if err != nil {
				return err
			}
			parser := language.NewParser(lang)
			parser.UnwrapProxyExports = cfg.UnwrapExports
			parser.Exclude = cfg.Exclude

			// TODO: perform the `explain` here
			println(fromFiles, toFiles)

			return nil
		},
	}

	return cmd
}
