package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/gabotechs/dep-tree/internal/config"
	"github.com/gabotechs/dep-tree/internal/js"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python"
	"github.com/gabotechs/dep-tree/internal/rust"
	"github.com/gabotechs/dep-tree/internal/utils"
)

var configPath string
var jsFollowTsConfigPaths bool
var followReExports bool
var exclude []string

var root *cobra.Command

func NewRoot(args []string) *cobra.Command {
	if args == nil {
		args = os.Args[1:]
	}

	root = &cobra.Command{
		Use:          "dep-tree",
		Version:      "v0.15.0",
		Short:        "Visualize and check your project's dependency tree",
		SilenceUsage: true,
		Args:         cobra.ArbitraryArgs,
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
	root.AddCommand(ConfigCmd())

	root.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to dep-tree's config file (default .dep-tree.yml)")
	root.PersistentFlags().BoolVar(&followReExports, "follow-re-exports", true, "whether to follow re-exports or not while resolving imports between files")
	root.PersistentFlags().BoolVar(&jsFollowTsConfigPaths, "js-follow-ts-config-paths", false, "whether to follow the tsconfig.json paths while resolving imports or not (default false)")
	root.PersistentFlags().StringArrayVar(&exclude, "exclude", nil, "Files that match this glob pattern will be ignored. You can provide an arbitrary number of --exclude flags")

	switch {
	case len(args) > 0 && utils.InArray(args[0], []string{"help", "completion", "-v", "--version", "-h", "--help"}):
		// do nothing.
	case len(args) == 0:
		root.SetArgs([]string{"help"})
	default:
		cmd, _, err := root.Find(args)
		if err == nil && cmd.Use == root.Use && !errors.Is(cmd.Flags().Parse(args), pflag.ErrHelp) {
			root.SetArgs(append([]string{"render"}, args...))
		}
	}
	return root
}

func makeParserBuilder(entrypoint string, cfg *config.Config) (language.NodeParserBuilder, error) {
	switch {
	case utils.EndsWith(entrypoint, js.Extensions):
		return language.ParserBuilder(js.MakeJsLanguage, &cfg.Js, cfg), nil
	case utils.EndsWith(entrypoint, rust.Extensions):
		return language.ParserBuilder(rust.MakeRustLanguage, &cfg.Rust, cfg), nil
	case utils.EndsWith(entrypoint, python.Extensions):
		return language.ParserBuilder(python.MakePythonLanguage, &cfg.Python, cfg), nil
	default:
		return nil, fmt.Errorf("file \"%s\" not supported", entrypoint)
	}
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.ParseConfig(configPath)
	if root.PersistentFlags().Changed("follow-re-exports") {
		cfg.FollowReExports = &followReExports
	}
	if root.PersistentFlags().Changed("js-follow-ts-config-paths") {
		cfg.Js.FollowTsConfigPaths = jsFollowTsConfigPaths
	}
	cfg.Exclude = append(cfg.Exclude, exclude...)
	for _, exclusion := range cfg.Exclude {
		if _, err := utils.GlobstarMatch(exclusion, ""); err != nil {
			return nil, fmt.Errorf("exclude pattern '%s' is not correctly formatted", exclusion)
		}
	}
	// Config load fails if a path was explicitly specified but the path does not exist.
	// If a path was not specified it's fine even if the default path does not exist.
	if os.IsNotExist(err) {
		if configPath != "" {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return cfg, nil
}
