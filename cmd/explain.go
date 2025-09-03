package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	regv1 "github.com/naivary/codemark/api/registry/v1"
	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/internal/console"
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
	domain := optionutil.DomainOf(ident)
	gen, err := mngr.Get(domain)
	if err != nil {
		return err
	}
	if optionutil.IsFQIdent(ident) {
		return e.explainOpt(ident, gen.Registry())
	}
	return nil
}

func (e *explainCmd) explainOpt(ident string, reg regv1.Registry) error {
	doc, err := reg.DocOf(ident)
	if err != nil {
		return err
	}
	fmt.Println(doc)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if doc.Default == "" {
		doc.Default = "<none>"
	}
	fmt.Fprintf(w, "DEFAULT: %s\n", doc.Default)
	fmt.Fprintf(w, "TYPE: <%s>\n", doc.Type)
	fmt.Println("DESC:")
	trunced := console.Trunc(doc.Desc, 70)
	for line := range strings.SplitSeq(trunced, "\n") {
		fmt.Fprintf(w, "\t\t%s\n", line)
	}
	return w.Flush()
}
