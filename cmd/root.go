package cmd

import (
	"github.com/spf13/cobra"
)

var Root = &cobra.Command{
	Use:          "dep-tree",
	Short:        "Visualize and check your project's dependency tree",
	SilenceUsage: true,
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
