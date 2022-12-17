package cmd

import (
	"errors"
	"fmt"
	"strings"

	"dep-tree/internal/graph"
	"dep-tree/internal/js"

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
		entrypoint := args[0]

		if endsWith(entrypoint, js.Extensions) {
			parser, err := js.MakeJsParser(entrypoint)
			if err != nil {
				return err
			}
			content, err := graph.RenderGraph[js.Data](entrypoint, parser)
			fmt.Print(content)
			return err
		} else {
			return errors.New("file not supported")
		}
	},
}
