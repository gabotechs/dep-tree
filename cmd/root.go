package cmd

import (
	"errors"
	"strings"

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
			parser, err := js.MakeJsParser(entrypoint)
			if err != nil {
				return err
			}

			return tui.Loop[js.Data](ctx, entrypoint, parser, nil)
		} else {
			return errors.New("file not supported")
		}
	},
}
