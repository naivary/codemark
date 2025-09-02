package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/optionutil"
	"github.com/naivary/codemark/outputer"
)

type explainCmd struct{}

func makeExplainCmd(genMngr *generator.Manager, outMngr *outputer.Manager) *cobra.Command {
	e := &explainCmd{}
	cmd := &cobra.Command{
		Use:   "explain",
		Short: "generate the artifacts for the given pattern",
		RunE:  e.runE(genMngr, outMngr),
		Args:  cobra.ExactArgs(2),
	}
	return cmd
}

func (e *explainCmd) runE(genMngr *generator.Manager, outMngr *outputer.Manager) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		typ := args[0]
		name := args[1]
		switch typ {
		case "marker":
			return e.explainMarker(name, genMngr)
		case "outputer":
			return e.explainOutputer(name, outMngr)
		default:
			return errors.New("explanations can only be provided for outputer and marker")
		}
	}
}

func (e *explainCmd) explainMarker(ident string, genMngr *generator.Manager) error {
	domain := optionutil.DomainOf(ident)
	gen, err := genMngr.Get(domain)
	if err != nil {
		return err
	}
	explanation := gen.Explain(ident)
	fmt.Println(explanation)
	return nil
}

func (e *explainCmd) explainOutputer(name string, outMngr *outputer.Manager) error {
	outputer, err := outMngr.Get(name)
	if err != nil {
		return err
	}
	explanation := outputer.Explain()
	fmt.Println(explanation)
	return nil
}
