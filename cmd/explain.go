package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/internal/explain"
	"github.com/naivary/codemark/optionutil"
	"github.com/naivary/codemark/outputer"
)

// TODO: Explain is a custom function now really the resposnibilty of the
// generator to have a explainf unction. its more they provide the doc and we
// have to retrieve it correctly.
//
// codemark explain --kind outputer fs
// codemark explain --kind marker <domain>:<resource>

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
		if e.kind == "marker" {
			ident := args[0]
			return e.explainMarker(ident, genMngr)
		}
		return nil
	}
}

func (e *explainCmd) explainMarker(ident string, mngr *generator.Manager) error {
	gen, err := mngr.Get(optionutil.DomainOf(ident))
	if err != nil {
		return err
	}
	return explain.Ident(os.Stdout, gen, ident)
}
