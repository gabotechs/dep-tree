package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func rootHelper(cmd *cobra.Command, args []string) error {
	if len(args) == 0 || args[0] == "help" {
		_ = cmd.Help()
		os.Exit(0)
	}
	return nil
}

var configPath string

func NewRoot(args []string) *cobra.Command {
	if args == nil {
		args = os.Args[1:]
	}
	root := &cobra.Command{
		Use:          "dep-tree",
		Version:      "v0.14.0",
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
	root.SetArgs(args)

	root.AddCommand(RenderCmd())
	root.AddCommand(CheckCmd())

	root.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to dep-tree's config file")

	loadDefault(root, args)
	return root
}

func loadDefault(root *cobra.Command, args []string) {
	if len(args) > 0 {
		if args[0] == "help" || args[0] == "completion" {
			return
		}
	}
	cmd, _, err := root.Find(args)
	if err == nil && cmd.Use == root.Use && !errors.Is(cmd.Flags().Parse(args), pflag.ErrHelp) {
		args := append([]string{"render"}, args...)
		root.SetArgs(args)
	}
}
