package explain

import (
	"errors"
	"fmt"
	"io"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	genv1 "github.com/naivary/codemark/api/generator/v1"
	"github.com/naivary/codemark/optionutil"
)

const _truncLen = 75

func Ident(w io.Writer, gen genv1.Generator, ident string) error {
	if len(ident) == 0 {
		return errors.New("ident cannot be empty")
	}
	if optionutil.IsFQIdent(ident) {
		doc, err := gen.Registry().DocOf(ident)
		if err != nil {
			return err
		}
		return option(w, doc)
	}
	if resourceName := optionutil.ResourceOf(ident); resourceName != "" {
		res := resourceDocOf(gen.Resources(), resourceName)
		optDocs := optDocsOf(gen.Registry().All(), resourceName)
		return resource(w, res, optDocs)
	}
	return domain(w, gen.Domain(), gen.Resources())
}

func domain(w io.Writer, domain docv1.Domain, resources []docv1.Resource) error {
	const cols = "NAME\tDESCRIPTION\n"
	tw := newTabWriter(w)
	fmt.Fprintf(tw, "%s\n\n", trunc(domain.Desc, _truncLen))
	fmt.Fprintf(tw, cols)
	for _, resource := range resources {
		desc := trunc(resource.Desc, _truncLen)
		name := fmt.Sprintf("%s:%s", domain.Name, resource.Name)
		writeLinesInCol(tw, "%s\t%s\n", desc, []any{name})
	}
	return tw.Flush()
}

func resource(w io.Writer, resource docv1.Resource, opts map[string]docv1.Option) error {
	const cols = "IDENT\tDEFAULT\tTYPE\tDESCRIPTION\n"
	// display resource descriptioon
	desc := trunc(resource.Desc, _truncLen)
	fmt.Println("DESCRIPTION")
	fmt.Printf("%s\n\n", desc)
	// display all options of the resource in a table

	tw := newTabWriter(w)
	fmt.Fprintf(tw, cols)
	for ident, doc := range opts {
		desc := trunc(doc.Desc, _truncLen)
		writeLinesInCol(tw, "%s\t%s\t%s\t%s\n", desc, []any{ident, doc.Default, doc.Type})
	}
	return tw.Flush()
}

func option(w io.Writer, opt *docv1.Option) error {
	if opt.Default == "" {
		opt.Default = "<none>"
	}
	if opt.Type == "" {
		opt.Type = "unknown"
	}
	if opt.Desc == "" {
		opt.Desc = "<none>"
	}
	tw := newTabWriter(w)
	fmt.Fprintf(tw, "DEFAULT: %s\n", opt.Default)
	fmt.Fprintf(tw, "TYPE: <%s>\n", opt.Type)
	fmt.Println("DESCRIPTION:")
	desc := trunc(opt.Desc, _truncLen)
	writeLinesInCol(w, " %s\n", desc, nil)
	return tw.Flush()
}
