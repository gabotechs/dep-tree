package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/gabotechs/dep-tree/internal/config"
	"github.com/gabotechs/dep-tree/internal/dummy"
	golang "github.com/gabotechs/dep-tree/internal/go"
	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/js"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/python"
	"github.com/gabotechs/dep-tree/internal/rust"
	"github.com/gabotechs/dep-tree/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const explainGroupId = "explain"
const renderGroupId = "render"
const checkGroupId = "check"
const defaultCommand = "entropy"

func NewRoot(args []string) *cobra.Command {
	if args == nil {
		args = os.Args[1:]
	}

	root := &cobra.Command{
		Use:               "dep-tree",
		Version:           "v0.22.2",
		Short:             "Visualize and check your project's dependency graph",
		SilenceUsage:      true,
		Args:              cobra.ArbitraryArgs,
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
		Example: `$ dep-tree src/index.ts
$ dep-tree entropy src/index.ts
$ dep-tree tree package/main.py --unwrap-exports
$ dep-tree check`,
		Long: `
      ____         _ __       _
     |  _ \   ___ |  _ \    _| |_  _ __  ___   ___
     | | | | / _ \| |_) |  |_   _||  __|/ _ \ / _ \
     | |_| ||  __/| .__/     | |  | |  |  __/|  __/
     |____/  \__| |_|        | \__|_|   \___| \___|
`,
	}

	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)
	root.SetArgs(args)

	root.AddGroup(&cobra.Group{ID: renderGroupId, Title: "Visualize your dependencies graphically"})
	root.AddGroup(&cobra.Group{ID: checkGroupId, Title: "Check your dependencies against your own rules"})
	root.AddGroup(&cobra.Group{ID: explainGroupId, Title: "Display what are the dependencies between two portions of code"})

	cliCfg := config.NewConfigCwd()

	var fileConfigPath string

	root.Flags().SortFlags = false
	root.PersistentFlags().SortFlags = false
	root.PersistentFlags().StringVarP(&fileConfigPath, "config", "c", "", "path to dep-tree's config file. (default .dep-tree.yml)")
	root.PersistentFlags().BoolVar(&cliCfg.UnwrapExports, "unwrap-exports", false, "trace re-exported symbols to the file where they are declared. (default false)")
	root.PersistentFlags().BoolVar(&cliCfg.Js.TsConfigPaths, "js-tsconfig-paths", true, "follow the tsconfig.json paths while resolving imports.")
	root.PersistentFlags().BoolVar(&cliCfg.Js.Workspaces, "js-workspaces", true, "take the workspaces attribute in the root package.json into account for resolving paths.")
	root.PersistentFlags().BoolVar(&cliCfg.Python.ExcludeConditionalImports, "python-exclude-conditional-imports", false, "exclude imports wrapped inside if or try statements. (default false)")
	root.PersistentFlags().StringArrayVar(&cliCfg.Only, "only", nil, "Files that do not match this glob pattern will be ignored. You can provide an arbitrary number of --only flags.")
	root.PersistentFlags().StringArrayVar(&cliCfg.Exclude, "exclude", nil, "Files that match this glob pattern will be ignored. You can provide an arbitrary number of --exclude flags.")

	cfgF := func() (*config.Config, error) {
		fileCfg, err := config.ParseConfigFromFile(fileConfigPath)
		if err != nil {
			return nil, err
		}
		fileCfg.EnsureAbsPaths()
		cliCfg.EnsureAbsPaths()

		// For these settings, prefer the CLI over the file config.
		for _, a := range []struct {
			name   string
			source *bool
			dest   *bool
		}{
			{"unwrap-exports", &cliCfg.UnwrapExports, &fileCfg.UnwrapExports},
			{"js-tsconfig-paths", &cliCfg.Js.TsConfigPaths, &fileCfg.Js.TsConfigPaths},
			{"js-workspaces", &cliCfg.Js.Workspaces, &fileCfg.Js.Workspaces},
			{"python-exclude-conditional-imports", &cliCfg.Python.ExcludeConditionalImports, &fileCfg.Python.ExcludeConditionalImports},
		} {
			if !root.PersistentFlags().Changed(a.name) {
				*a.dest = *a.source
			}
		}
		// NOTE: hard-enable this for now, as they don't produce a very good output.
		fileCfg.Python.IgnoreFromImportsAsExports = true
		fileCfg.Python.IgnoreDirectoryImports = true

		// merge the exclusions and inclusions from the CLI and from the config.
		fileCfg.Exclude = append(fileCfg.Exclude, cliCfg.Exclude...)
		fileCfg.Only = append(fileCfg.Only, cliCfg.Only...)

		return fileCfg, fileCfg.ValidatePatterns()
	}

	root.AddCommand(
		EntropyCmd(cfgF),
		TreeCmd(cfgF),
		CheckCmd(cfgF),
		ConfigCmd(cfgF),
		ExplainCmd(cfgF),
	)

	switch {
	case len(args) > 0 && utils.InArray(args[0], []string{"help", "completion", "-v", "--version", "-h", "--help"}):
		// do nothing.
	case len(args) == 0:
		// if not args where provided, default to help.
		root.SetArgs([]string{"help"})
	default:
		// if some args where provided, but it's none of the root commands,
		// choose a default command.
		cmd, _, err := root.Find(args)
		if err == nil && cmd.Use == root.Use && !errors.Is(cmd.Flags().Parse(args), pflag.ErrHelp) {
			root.SetArgs(append([]string{defaultCommand}, args...))
		}
	}
	return root
}

//nolint:gocyclo
func inferLang(files []string, cfg *config.Config) (language.Language, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("at least 1 file must be provided for infering the language")
	}
	score := struct {
		js     int
		python int
		rust   int
		golang int
		dummy  int
	}{}
	top := struct {
		lang string
		v    int
	}{}
	for _, file := range files {
		switch {
		case utils.EndsWith(file, js.Extensions):
			score.js += 1
			if score.js > top.v {
				top.v = score.js
				top.lang = "js"
			}
		case utils.EndsWith(file, rust.Extensions):
			score.rust += 1
			if score.rust > top.v {
				top.v = score.rust
				top.lang = "rust"
			}
		case utils.EndsWith(file, python.Extensions):
			score.python += 1
			if score.python > top.v {
				top.v = score.python
				top.lang = "python"
			}
		case utils.EndsWith(file, golang.Extensions):
			score.golang += 1
			if score.golang > top.v {
				top.v = score.golang
				top.lang = "golang"
			}
		case utils.EndsWith(file, dummy.Extensions):
			score.dummy += 1
			if score.dummy > top.v {
				top.v = score.dummy
				top.lang = "dummy"
			}
		}
	}
	if top.lang == "" {
		return nil, errors.New("none of the provided files belong to the a supported language")
	}
	switch top.lang {
	case "js":
		return js.MakeJsLanguage(&cfg.Js)
	case "rust":
		return rust.MakeRustLanguage(&cfg.Rust)
	case "python":
		return python.MakePythonLanguage(&cfg.Python)
	case "golang":
		return golang.NewLanguage(files[0], &cfg.Golang)
	case "dummy":
		return &dummy.Language{}, nil
	default:
		return nil, fmt.Errorf("file \"%s\" not supported", files[0])
	}
}

func filesFromArgs(args []string) ([]string, error) {
	var result []string
	for _, arg := range args {
		basepath, pattern := doublestar.SplitPattern(arg)
		fsys := os.DirFS(basepath)
		matches, err := doublestar.Glob(fsys, pattern)
		if err != nil {
			return nil, err
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("%s does not match with any existing file", arg)
		}
		for _, match := range matches {
			abs, err := filepath.Abs(filepath.Join(basepath, match))
			if err != nil {
				return nil, err
			} else if !utils.FileExists(abs) {
				if !utils.DirExists(abs) {
					return nil, fmt.Errorf("file %s does not exist", match)
				}
			} else {
				result = append(result, abs)
			}
		}
	}
	if len(result) == 0 {
		return result, errors.New("no valid files where provided")
	}

	return result, nil
}

func applyConfigToParser(parser *language.Parser, cfg *config.Config) {
	parser.UnwrapProxyExports = cfg.UnwrapExports
	parser.Exclude = cfg.Exclude
	parser.Include = cfg.Only
}

func relPathDisplay(node *graph.Node[*language.FileInfo]) string {
	return node.Data.RelPath
}
