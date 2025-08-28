package cmd

import (
	"github.com/spf13/cobra"
)

type explainCmd struct{}

func makeExplainCmd() *cobra.Command {
	e := &explainCmd{}
	cmd := &cobra.Command{
		Use:  "explain [ident]",
		Args: cobra.ExactArgs(1),
		RunE: e.runE,
	}
	return cmd
}

func (e *explainCmd) runE(cmd *cobra.Command, args []string) error {
	return nil
}
