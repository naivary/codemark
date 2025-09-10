package explain

import (
	"fmt"
	"io"
	"strings"

	outv1 "github.com/naivary/codemark/api/outputer/v1"
)

func Outputer(w io.Writer, out outv1.Outputer) error {
	doc := out.Doc()
	if doc.Desc == "" {
		doc.Desc = _none
	}
	fmt.Fprintf(w, "NAME: %s\n", doc.Name)
	fmt.Fprintf(w, "DESCRIPTION: %s\n", trunc(doc.Desc, _truncLen))
	fmt.Println("FLAGS:")
	usage := strings.TrimSpace(out.Flags().FlagUsages())
	if usage == "" {
		usage = _none
	}
	fmt.Fprintf(w, "  %s\n", usage)
	return nil
}

func AllOutputer(w io.Writer, outs []outv1.Outputer) error {
	const cols = "NAME\tDESCRIPTION\n"
	tw := newTabWriter(w)
	fmt.Fprintf(tw, cols)
	for _, out := range outs {
		doc := out.Doc()
		desc := trunc(doc.Desc, _truncLen)
		writeLinesInCol(tw, "%s\t%s\n", desc, []any{doc.Name})
	}
	return tw.Flush()
}
