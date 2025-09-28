package cmd

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/internal/explain"
	"github.com/naivary/codemark/optionutil"
	"github.com/naivary/codemark/outputer"
)

type explainCmd struct {
	kind string
}

func makeExplainCmd(genMngr *generator.Manager, outMngr *outputer.Manager) *cobra.Command {
	e := &explainCmd{}
	cmd := &cobra.Command{
		Use:   "explain",
		Short: "get the documentation about outputer or marker",
		RunE:  e.runE(genMngr, outMngr),
		Args:  cobra.ExactArgs(1),
	}
	cmd.Flags().StringVar(&e.kind, "kind", "gen", "kind of explanation")
	return cmd
}

func (e *explainCmd) runE(genMngr *generator.Manager, outMngr *outputer.Manager) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		name := args[0]
		switch e.kind {
		case "generator", "gen":
			return e.explainGen(name, genMngr)
		case "outputer":
			return e.explainOutputer(name, outMngr)
		case "config":
			return e.explainConfig(name, genMngr)
		default:
			return fmt.Errorf("unknown kind: %s", e.kind)
		}
	}
}

func (e *explainCmd) explainGen(ident string, mngr *generator.Manager) error {
	if ident == "all" {
		gens := slices.Collect(maps.Values(mngr.Gens()))
		return explain.AllGens(os.Stdout, gens)
	}
	gen, err := mngr.Get(optionutil.DomainOf(ident))
	if err != nil {
		return err
	}
	return explain.Generator(os.Stdout, gen, ident)
}

func (e *explainCmd) explainOutputer(name string, mngr *outputer.Manager) error {
	if name == "all" {
		all := slices.Collect(maps.Values(mngr.All()))
		return explain.AllOutputer(os.Stdout, all)
	}
	out, err := mngr.Get(name)
	if err != nil {
		return err
	}
	return explain.Outputer(os.Stdout, out)
}

func (e *explainCmd) explainConfig(dotPath string, mngr *generator.Manager) error {
	paths := strings.Split(dotPath, ".")
	name := paths[0]
	gen, err := mngr.Get(name)
	if err != nil {
		return err
	}
	return explain.Config(os.Stdout, strings.Join(paths[1:], "."), gen.ConfigDoc())
}
