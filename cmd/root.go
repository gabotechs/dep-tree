package cmd

import (
	"strings"

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
	Use:   "dep-tree",
	Short: "Visualize and check your project's dependency tree",
	Long: `
      ____         _ __       _
     |  _ \   ___ |  _ \    _| |_  _ __  ___   ___ 
     | | | | / _ \| |_) |  |_   _||  __|/ _ \ / _ \
     | |_| ||  __/| .__/     | |  | |  |  __/|  __/
     |____/  \__| |_|        | \__|_|   \___| \___|
`,
}

func init() {
	Root.AddCommand(RenderCmd())
	Root.AddCommand(CheckCmd())
}
