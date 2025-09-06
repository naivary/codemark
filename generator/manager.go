package generator

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
	"text/tabwriter"

	convv1 "github.com/naivary/codemark/api/converter/v1"
	docv1 "github.com/naivary/codemark/api/doc/v1"
	genv1 "github.com/naivary/codemark/api/generator/v1"
	regv1 "github.com/naivary/codemark/api/registry/v1"
	"github.com/naivary/codemark/internal/config"
	"github.com/naivary/codemark/internal/console"
	"github.com/naivary/codemark/loader"
	"github.com/naivary/codemark/optionutil"
	"github.com/naivary/codemark/registry"
)

type domain = string

type Manager struct {
	gens map[domain]genv1.Generator

	cfg map[string]any
}

func NewManager(cfgFile string, gens ...genv1.Generator) (*Manager, error) {
	const configSection = "gens"
	mngr := &Manager{
		gens: make(map[domain]genv1.Generator),
	}
	cfg, err := config.ReadIn(cfgFile, configSection)
	if err != nil {
		return nil, err
	}
	mngr.cfg = cfg
	for _, gen := range gens {
		err := mngr.Add(gen)
		if err != nil {
			return nil, err
		}
	}
	return mngr, nil
}

func (m *Manager) Domains() []string {
	return slices.Collect(maps.Keys(m.gens))
}

func (m *Manager) All() map[domain]genv1.Generator {
	return m.gens
}

func (m *Manager) Generate(convs []convv1.Converter, pattern string) (map[domain][]*genv1.Artifact, error) {
	reg, err := m.merge(m.allGens())
	if err != nil {
		return nil, err
	}
	info, err := loader.Load(reg, convs, pattern)
	if err != nil {
		return nil, err
	}
	output := make(map[domain][]*genv1.Artifact)
	for _, gen := range m.gens {
		artifacts, err := gen.Generate(info, m.configFor(gen))
		if err != nil {
			return nil, err
		}
		output[gen.Domain().Name] = artifacts
	}
	return output, nil
}

func (m *Manager) Get(domain string) (genv1.Generator, error) {
	gen, found := m.gens[domain]
	if !found {
		return nil, fmt.Errorf("generator not found for domain: %s", domain)
	}
	return gen, nil
}

func (m *Manager) Add(gen genv1.Generator) error {
	domain := gen.Domain().Name
	if _, found := m.gens[domain]; found {
		return fmt.Errorf("generator for domain already exists: %s", domain)
	}
	m.gens[domain] = gen
	return nil
}

func (m *Manager) Explain(w io.Writer, ident string) error {
	if optionutil.OptionOf(ident) != "" {
		return m.explainOption(w, ident)
	}
	if optionutil.ResourceOf(ident) != "" {
		return m.explainResource(w, ident)
	}
	if optionutil.DomainOf(ident) != "" {
		return m.explainDomain(w, ident)
	}
	return fmt.Errorf("no explanation could be found for `%s`. Make sure you used the correct syntax of <domain>:<resource>:<option>")
}

func (m *Manager) explainDomain(w io.Writer, ident string) error {
	return nil
}

func (m *Manager) explainResource(w io.Writer, ident string) error {
	domain := optionutil.DomainOf(ident)
	resourceName := optionutil.ResourceOf(ident)
	gen, err := m.Get(domain)
	if err != nil {
		return err
	}
	var resource docv1.Resource
	for _, res := range gen.Resources() {
		if resource.Name == resourceName {
			resource = res
			break
		}
	}
	optionsOfResource := make(map[string]*docv1.Option, 0)
	for fqi, opt := range gen.Registry().All() {
		if optionutil.ResourceOf(fqi) == resourceName {
			optionsOfResource[fqi] = opt.Doc
		}
	}
	tw := m.newTabWriter(w)
	fmt.Fprintln(tw, "IDENT\tDEFAULT\tTYPE\tDESC")
	for fqi, optDoc := range optionsOfResource {
		desc := console.Trunc(optDoc.Desc, 75)
		lines := strings.Split(desc, "\n")
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", fqi, optDoc.Default, optDoc.Type, lines[0])
		for _, line := range lines[1:] {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", "", "", "", line)
		}
	}

	// fmt.Fprintf(tw, "%s\n\n", resource.Desc)
	// fmt.Println("OPTIONS:")
	// const lineLen = 75
	// for fqi, optDoc := range optionsOfResource {
	// 	fmt.Fprintf(tw, "  IDENT: %s\n", fqi)
	// 	fmt.Fprintf(tw, "  DEFAULT: %s\n", optDoc.Default)
	// 	fmt.Fprintf(tw, "  TYPE: <%s>\n", optDoc.Type)
	// 	fmt.Fprintf(tw, "  DESCRIPTION:\n")
	// 	for line := range strings.SplitSeq(console.Trunc(optDoc.Desc, lineLen), "\n") {
	// 		fmt.Fprintf(tw, "    %s\n", line)
	// 	}
	// 	fmt.Fprintf(tw, "-------\n")
	// }
	return tw.Flush()
}

func (m *Manager) explainOption(w io.Writer, ident string) error {
	domain := optionutil.DomainOf(ident)
	gen, err := m.Get(domain)
	if err != nil {
		return err
	}
	doc, err := gen.Registry().DocOf(ident)
	if err != nil {
		return err
	}
	const lineLen = 75
	tw := m.newTabWriter(w)
	fmt.Fprintf(tw, "DEFAULT: %s\n", doc.Default)
	fmt.Fprintf(tw, "TYPE: <%s>\n", doc.Type)
	fmt.Fprintf(tw, "DESCRIPTION: %s\n", console.Trunc(doc.Desc, lineLen))
	return tw.Flush()
}

func (m *Manager) newTabWriter(w io.Writer) *tabwriter.Writer {
	const (
		minWidth = 0
		tabWidth = 0
		padding  = 2
		padChar  = ' '
		flags    = 0
	)
	return tabwriter.NewWriter(w, minWidth, tabWidth, padding, padChar, flags)
}

func (m *Manager) allGens() []genv1.Generator {
	return slices.Collect(maps.Values(m.gens))
}

func (m *Manager) configFor(gen genv1.Generator) map[string]any {
	genCfg, isMap := m.cfg[gen.Domain().Name].(map[string]any)
	if isMap {
		return genCfg
	}
	return make(map[string]any)
}

func (m *Manager) merge(gens []genv1.Generator) (regv1.Registry, error) {
	regs := make([]regv1.Registry, 0, len(gens))
	for _, gen := range gens {
		regs = append(regs, gen.Registry())
	}
	return registry.Merge(regs...)
}
