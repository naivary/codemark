package explainer

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	genv1 "github.com/naivary/codemark/api/generator/v1"
	"github.com/naivary/codemark/generator"
	"github.com/naivary/codemark/internal/console"
	"github.com/naivary/codemark/optionutil"
)

var _ Explainer = (*markerExplainer)(nil)

type markerExplainer struct {
	mngr     *generator.Manager
	truncLen int
}

func NewMarkerExplainer(mngr *generator.Manager) Explainer {
	return &markerExplainer{
		mngr:     mngr,
		truncLen: 75,
	}
}

func (m markerExplainer) newTabWriter(w io.Writer) *tabwriter.Writer {
	const (
		minWidth = 0
		tabWidth = 0
		padding  = 2
		padChar  = ' '
		flags    = 0
	)
	return tabwriter.NewWriter(w, minWidth, tabWidth, padding, padChar, flags)
}

func (m markerExplainer) Explain(w io.Writer, args ...string) error {
	ident := args[0]
	domain := optionutil.DomainOf(ident)
	gen, err := m.mngr.Get(domain)
	if err != nil {
		return err
	}
	if optionutil.OptionOf(ident) != "" {
		return m.explainOption(w, gen, ident)
	}
	if resourceName := optionutil.ResourceOf(ident); resourceName != "" {
		return m.explainResource(w, gen, resourceName)
	}
	return m.explainDomain(w, gen)
}

func (m markerExplainer) explainDomain(w io.Writer, gen genv1.Generator) error {
	resources := gen.Resources()
	tw := m.newTabWriter(w)
	fmt.Fprintf(tw, "%s\n\n", console.Trunc(gen.Domain().Desc, m.truncLen))
	fmt.Fprintln(tw, "NAME\tDESCRIPTION")
	for _, resource := range resources {
		desc := console.Trunc(resource.Desc, m.truncLen)
		lines := strings.Split(desc, "\n")
		name := fmt.Sprintf("%s:%s", gen.Domain().Name, resource.Name)
		fmt.Fprintf(tw, "%s\t%s\t\n", name, lines[0])
		for _, line := range lines[1:] {
			fmt.Fprintf(tw, "%s\t%s\n", "", line)
		}
	}
	return tw.Flush()
}

func (m markerExplainer) explainResource(w io.Writer, gen genv1.Generator, resourceName string) error {
	for _, resource := range gen.Resources() {
		if resource.Name != resourceName {
			continue
		}
		desc := console.Trunc(resource.Desc, m.truncLen)
		fmt.Println("DESCRIPTION")
		fmt.Printf("%s\n\n", desc)
	}
	optionsOfResource := make(map[string]*docv1.Option, 0)
	for fqi, opt := range gen.Registry().All() {
		if optionutil.ResourceOf(fqi) == resourceName {
			optionsOfResource[fqi] = opt.Doc
		}
	}
	tw := m.newTabWriter(w)
	fmt.Fprintln(tw, "IDENT\tDEFAULT\tTYPE\tDESCRIPTION")
	for fqi, optDoc := range optionsOfResource {
		prepareOptDoc(optDoc)
		desc := console.Trunc(optDoc.Desc, m.truncLen)
		lines := strings.Split(desc, "\n")
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", fqi, optDoc.Default, optDoc.Type, lines[0])
		for _, line := range lines[1:] {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", "", "", "", line)
		}
	}
	return tw.Flush()
}

func (m markerExplainer) explainOption(w io.Writer, gen genv1.Generator, ident string) error {
	domain := optionutil.DomainOf(ident)
	gen, err := m.mngr.Get(domain)
	if err != nil {
		return err
	}
	doc, err := gen.Registry().DocOf(ident)
	if err != nil {
		return err
	}
	prepareOptDoc(doc)
	tw := m.newTabWriter(w)
	fmt.Fprintf(tw, "DEFAULT: %s\n", doc.Default)
	fmt.Fprintf(tw, "TYPE: <%s>\n", doc.Type)
	fmt.Println("DESCRIPTION:")
	for line := range strings.SplitSeq(console.Trunc(doc.Desc, m.truncLen), "\n") {
		fmt.Fprintf(tw, "  %s\n", line)
	}
	return tw.Flush()
}

func prepareOptDoc(optDoc *docv1.Option) {
	if optDoc.Default == "" {
		optDoc.Default = "<none>"
	}
	if optDoc.Type == "" {
		optDoc.Type = "unknown"
	}
	if optDoc.Desc == "" {
		optDoc.Desc = "<none>"
	}

}
