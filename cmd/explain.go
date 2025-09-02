package cmd

import (
	"github.com/spf13/cobra"

	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/outputer"
)

// TODO: Explain is a custom function now really the resposnibilty of the
// generator to have a explainf unction. its more they provide the doc and we
// have to retrieve it correctly.
//
// codemark explain --type outputer fs

type ExplainKind string

var (
	ExplainKindMarker   ExplainKind = "marker"
	ExplainKindOutputer ExplainKind = "outputer"
)

type explainCmd struct {
	kind ExplainKind
}

func makeExplainCmd(genMngr *generator.Manager, outMngr *outputer.Manager) *cobra.Command {
	e := &explainCmd{}
	cmd := &cobra.Command{
		Use:   "explain",
		Short: "get the documentation about outputer or marker",
		RunE:  e.runE(genMngr, outMngr),
		Args:  cobra.ExactArgs(2),
	}
	return cmd
}

func (e *explainCmd) runE(genMngr *generator.Manager, outMngr *outputer.Manager) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}
