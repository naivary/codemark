package cmd

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/option"
)

func makeExplainCmd(mngr *generator.Manager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "explain [ident]",
		Short: "",
		Long:  "",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			ident := args[len(args)-1]
			identParts := strings.Split(ident, ":")
			if len(identParts) != 3 {
				return fmt.Errorf("full qualified identifier needed: %s", ident)
			}
			if slices.Contains(identParts, "") {
				return errors.New("all parts of the identifier have to have at least one letter")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ident := args[len(args)-1]
			domain := option.DomainOf(ident)
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
