package explain

import (
	"errors"
	"fmt"
	"io"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	genv1 "github.com/naivary/codemark/api/generator/v1"
	optv1 "github.com/naivary/codemark/api/option/v1"
	"github.com/naivary/codemark/optionutil"
)

const _truncLen = 75

func Ident(w io.Writer, gen genv1.Generator, ident string) error {
	if len(ident) == 0 {
		return errors.New("ident cannot be empty")
	}
	if optionutil.IsFQIdent(ident) {
		opt, err := gen.Registry().Get(ident)
		if err != nil {
			return err
		}
		return option(w, opt)
	}
	if resourceName := optionutil.ResourceOf(ident); resourceName != "" {
		resc := gen.Resources()[resourceName]
		if resc == nil {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		opts := optsOf(gen.Registry().All(), resourceName)
		return resource(w, *resc, opts)
	}
	return domain(w, gen.Domain(), gen.Resources())
}

func domain(w io.Writer, domain docv1.Domain, resources map[string]*docv1.Resource) error {
	const cols = "NAME\tDESCRIPTION\n"
	tw := newTabWriter(w)
	fmt.Fprintf(tw, "%s\n\n", trunc(domain.Desc, _truncLen))
	fmt.Fprintf(tw, cols)
	for name, resource := range resources {
		desc := trunc(resource.Desc, _truncLen)
		name := fmt.Sprintf("%s:%s", domain.Name, name)
		writeLinesInCol(tw, "%s\t%s\n", desc, []any{name})
	}
	return tw.Flush()
}

func resource(w io.Writer, resource docv1.Resource, opts map[string]*optv1.Option) error {
	const cols = "IDENT\tDEFAULT\tTYPE\tTARGETS\tDESCRIPTION\n"
	// display resource descriptioon
	desc := trunc(resource.Desc, _truncLen)
	fmt.Println("DESCRIPTION")
	fmt.Printf("%s\n\n", desc)
	// display all options of the resource in a table
	tw := newTabWriter(w)
	fmt.Fprintf(tw, cols)
	for ident, opt := range opts {
		doc := opt.Doc
		desc := trunc(doc.Desc, _truncLen)
		writeLinesInCol(tw, "%s\t%s\t%s\t%s\t%s\n", desc, []any{ident, doc.Default, TypeOf(opt.Type), targetsToString(opt.Targets)})
	}
	return tw.Flush()
}

func option(w io.Writer, opt *optv1.Option) error {
	doc := opt.Doc
	if doc.Default == "" {
		doc.Default = _none
	}
	if doc.Desc == "" {
		doc.Desc = _none
	}
	tw := newTabWriter(w)
	fmt.Fprintf(tw, "DEFAULT: %s\n", doc.Default)
	fmt.Fprintf(tw, "TYPE: %s\n", TypeOf(opt.Type))
	fmt.Fprintf(tw, "TARGETS: %s\n", targetsToString(opt.Targets))
	fmt.Println("DESCRIPTION:")
	desc := trunc(doc.Desc, _truncLen)
	writeLinesInCol(w, "   %s\n", desc, nil)
	return tw.Flush()
}
