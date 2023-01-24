package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
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
				_, dt, err := dep_tree.NewDepTree[js.Data](ctx, parser)
				if err != nil {
					return err
				}
				return dt.Validate(cfg, func(nodeId string) string {
					processed, err := filepath.Rel(cfg.Path, nodeId)
					if err != nil {
						return nodeId
					}
					return processed
				})
			} else {
				return tui.Loop[js.Data](
					ctx,
					entrypoint,
					js.MakeJsParser,
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
