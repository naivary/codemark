package cmd

import (
	"fmt"
	"maps"
	"os"
	"slices"

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
	cmd.Flags().StringVar(&e.kind, "kind", "marker", "kind of explanation")
	return cmd
}

func (e *explainCmd) runE(genMngr *generator.Manager, outMngr *outputer.Manager) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		name := args[0]
		switch e.kind {
		case "marker":
			return e.explainMarker(name, genMngr)
		case "outputer":
			return e.explainOutputer(name, outMngr)
		case "all":
			return e.explainOutputer(name, outMngr)
		default:
			return fmt.Errorf("unknown kind: %s", e.kind)
		}
	}
}

func (e *explainCmd) explainMarker(ident string, mngr *generator.Manager) error {
	gen, err := mngr.Get(optionutil.DomainOf(ident))
	if err != nil {
		return err
	}
	return explain.Ident(os.Stdout, gen, ident)
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
