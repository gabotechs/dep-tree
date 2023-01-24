package cmd

import (
	"errors"
	"fmt"
	"strings"

	"dep-tree/internal/config"
	"dep-tree/internal/dep_tree"
	"dep-tree/internal/js"
	"dep-tree/internal/tui"

	"github.com/spf13/cobra"
)

func endsWith(string string, substrings []string) bool {
	for _, substring := range substrings {
		if strings.HasSuffix(string, substring) {
			return true
		}
	}
	return false
}

var check bool
var configPath string

var Root = &cobra.Command{
	Use:   "<path>",
	Short: "Render your project's dependency tree",
	Long: `
      ____         _ __       _
     |  _ \   ___ |  _ \    _| |_  _ __  ___   ___ 
     | | | | / _ \| |_) |  |_   _||  __|/ _ \ / _ \
     | |_| ||  __/| .__/     | |  | |  |  __/|  __/
     |____/  \__| |_|        | \__|_|   \___| \___|

`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		entrypoint := args[0]

		if endsWith(entrypoint, js.Extensions) {

			if check {
				cfg, err := config.ParseConfig(configPath)
				if err != nil {
					return fmt.Errorf("could not parse config file %s: %w", configPath, err)
				}
				parser, err := js.MakeJsParser(entrypoint)
				if err != nil {
					return err
				}
				ctx, dt, err := dep_tree.NewDepTree[js.Data](ctx, parser)
				if err != nil {
					return err
				}
				_, err = dt.Validate(ctx, cfg)
				return err
			} else {
				return tui.Loop[js.Data](
					ctx,
					entrypoint,
					// NOTE: it should be sufficient to pass js.MakeJsParser, but go complains.
					func(s string) (dep_tree.NodeParser[js.Data], error) {
						return js.MakeJsParser(s)
					},
					nil,
				)
			}

		} else {
			return errors.New("file not supported")
		}
	},
}

func init() {
	Root.Flags().BoolVar(&check, "check", false, "check if the dependency graph matches the rules defined by the user")
	Root.Flags().StringVar(&configPath, "config", ".dep-tree.yml", "path to dep-tree's config file")
}
