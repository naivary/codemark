package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/optionutil"
)

func makeExplainCmd(mngr *generator.Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "explain [ident]",
		Short: "",
		Long:  "",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ident := args[len(args)-1]
			domain := optionutil.DomainOf(ident)
			if domain == "" {
				domain = ident
			}
			gen, err := mngr.Get(domain)
			if err != nil {
				return err
			}
			explanation := gen.Explain(ident)
			if explanation == "" {
				return fmt.Errorf("no explanation found for %s", ident)
			}
			_, err = fmt.Fprintln(os.Stdout, explanation)
			return err
		},
	}
	return cmd
}
