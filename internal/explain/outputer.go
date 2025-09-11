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
	fmt.Printf("FLAGS:")
	usage := _none
	// inline format if usage is none
	format := " %s\n"
	flags := out.Flags()
	if flags != nil {
		format = "\n  %s\n"
		usage = strings.TrimSpace(out.Flags().FlagUsages())
	}
	fmt.Fprintf(w, format, usage)
	return nil
}

func AllOutputer(w io.Writer, outs []outv1.Outputer) error {
	const cols = "NAME\tSUMMARY\n"
	tw := newTabWriter(w)
	fmt.Fprintf(tw, cols)
	for _, out := range outs {
		doc := out.Doc()
		summary := trunc(doc.Summary, _truncLen)
		writeLinesInCol(tw, "%s\t%s\n", summary, []any{doc.Name})
	}
	return tw.Flush()
}
