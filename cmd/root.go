package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func rootHelper(cmd *cobra.Command, args []string) error {
	if len(args) == 0 || args[0] == "help" {
		_ = cmd.Help()
		os.Exit(0)
	}
	return nil
}

func NewRoot() *cobra.Command {
	root := &cobra.Command{
		Use:          "dep-tree",
		Version:      "v0.13.7",
		Short:        "Visualize and check your project's dependency tree",
		SilenceUsage: true,
		Args:         rootHelper,
		RunE:         runRender,
		Long: `
      ____         _ __       _
     |  _ \   ___ |  _ \    _| |_  _ __  ___   ___
     | | | | / _ \| |_) |  |_   _||  __|/ _ \ / _ \
     | |_| ||  __/| .__/     | |  | |  |  __/|  __/
     |____/  \__| |_|        | \__|_|   \___| \___|
`,
	}

	root.AddCommand(RenderCmd())
	root.AddCommand(CheckCmd())

	return root
}
