package cmd

import (
	"errors"
	"os"
	"slices"
	"strings"

	"github.com/gabotechs/dep-tree/internal/config"
	"github.com/gabotechs/dep-tree/internal/explain"
	"github.com/gabotechs/dep-tree/internal/graph"
	"github.com/gabotechs/dep-tree/internal/language"
	"github.com/gabotechs/dep-tree/internal/utils"
	"github.com/spf13/cobra"
)

func ExplainCmd(cfgF func() (*config.Config, error)) *cobra.Command {
	var overlapLeft bool
	var overlapRight bool

	cmd := &cobra.Command{
		Use:     "explain",
		Short:   "Shows all the dependencies between two parts of the code",
		GroupID: explainGroupId,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if overlapLeft && overlapRight {
				return errors.New("only one of --overlap-left (-l) or --overlap-right (-r) can be used at a time")
			}

			fromFiles, err := filesFromArgs([]string{args[0]})
			if err != nil {
				return err
			}

			toFiles, err := filesFromArgs([]string{args[1]})
			if err != nil {
				return err
			}

			if overlapLeft {
				toFiles = utils.RemoveOverlap(toFiles, fromFiles)
			}

			if overlapRight {
				fromFiles = utils.RemoveOverlap(fromFiles, toFiles)
			}

			cfg, err := cfgF()
			if err != nil {
				return err
			}

			lang, err := inferLang(fromFiles, cfg)
			if err != nil {
				return err
			}

			parser := language.NewParser(lang)
			applyConfigToParser(parser, cfg)

			cwd, _ := os.Getwd()
			tempCfg := config.Config{Path: cwd, Only: args}
			tempCfg.EnsureAbsPaths()
			parser.Include = tempCfg.Only

			deps, err := explain.Explain[*language.FileInfo](
				parser,
				fromFiles,
				toFiles,
				graph.NewStdErrCallbacks[*language.FileInfo](relPathDisplay),
			)
			if err != nil {
				return err
			}

			// If more than 1 package is referenced, display it.
			shouldIncludePackagePrefix := moreThanOnePackage(deps)

			rendered := make([]string, len(deps))
			for i, r := range deps {
				if shouldIncludePackagePrefix {
					fromPkg := r[0].Data.Package
					if strings.HasPrefix(fromPkg, "@") {
						fromPkg = fromPkg[1:]
					}
					toPkg := r[1].Data.Package
					if strings.HasPrefix(toPkg, "@") {
						toPkg = toPkg[1:]
					}
					rendered[i] = fromPkg + "@" + relPathDisplay(r[0]) + " -> " + toPkg + "@" + relPathDisplay(r[1])
				} else {
					rendered[i] = relPathDisplay(r[0]) + " -> " + relPathDisplay(r[1])
				}
			}

			slices.Sort(rendered)
			for _, line := range rendered {
				cmd.Println(line)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&overlapLeft, "overlap-left", "l", false, "When there's an overlap between the files at the left and the right, keep the ones at the left")
	cmd.Flags().BoolVarP(&overlapRight, "overlap-right", "r", false, "When there's an overlap between the files at the left and the right, keep the ones at the right")

	return cmd
}

func moreThanOnePackage(deps [][2]*graph.Node[*language.FileInfo]) bool {
	packages := map[string]struct{}{}
	for _, nodes := range deps {
		for _, node := range nodes {
			if _, ok := packages[node.Data.Package]; !ok {
				packages[node.Data.Package] = struct{}{}
				if len(packages) > 1 {
					return true
				}
			}
		}
	}
	return false
}
