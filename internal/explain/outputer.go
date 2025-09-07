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
