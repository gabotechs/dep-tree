package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/gabotechs/dep-tree/internal/config"
	"github.com/gabotechs/dep-tree/internal/utils"
)

var configPath string
var jsFollowTsConfigPaths bool
var cfg *config.Config

func loadConfig() error {
	var err error
	cfg, err = config.ParseConfig(configPath)
	if jsFollowTsConfigPaths {
		cfg.Js.FollowTsConfigPaths = true
	}
	if os.IsNotExist(err) {
		if configPath != "" {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func NewRoot(args []string) *cobra.Command {
	if args == nil {
		args = os.Args[1:]
	}
	root := &cobra.Command{
		Use:          "dep-tree",
		Version:      "v0.14.2",
		Short:        "Visualize and check your project's dependency tree",
		SilenceUsage: true,
		Args:         cobra.ArbitraryArgs,
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
	root.AddCommand(EntropyCmd())

	root.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to dep-tree's config file")
	root.PersistentFlags().BoolVar(&jsFollowTsConfigPaths, "js-follow-ts-config-paths", false, "whether to follow the tsconfig.json paths while resolving imports or not")

	loadDefault(root, args)
	return root
}

func loadDefault(root *cobra.Command, args []string) {
	if len(args) > 0 {
		if utils.InArray(args[0], []string{"help", "completion", "-v", "--version", "-h", "--help"}) {
			return
		}
	} else if len(args) == 0 {
		root.SetArgs([]string{"help"})
		return
	}
	cmd, _, err := root.Find(args)
	if err == nil && cmd.Use == root.Use && !errors.Is(cmd.Flags().Parse(args), pflag.ErrHelp) {
		args := append([]string{"render"}, args...)
		root.SetArgs(args)
	}
}
