package api

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/naivary/codemark/definition/target"
	"github.com/naivary/codemark/syntax"
)

// TODO: make this a custom package
func trunc(s string, n int) string {
	var b bytes.Buffer
	pos := 1
	for _, r := range s {
		if pos%n == 0 && r == ' ' {
			fmt.Fprintf(&b, "\n")
			pos = 1
			continue
		}
		if r == '\n' {
			pos = 1
		}
		if pos%n != 0 {
			pos++
		}
		fmt.Fprint(&b, string(r))
	}
	return b.String()
}

type OptionDoc struct {
	Targets []target.Target
	Ident   string
	Doc     string
	Default string
}

func (o OptionDoc) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "Ident: %s\n", o.Ident)
	fmt.Fprintf(&b, "Default: %s\n", o.Default)
	fmt.Fprintf(&b, "Targets: %v\n", o.Targets)
	fmt.Fprintf(&b, "%s", trunc(o.Doc, 80))
	return b.String()
}

type ResourceDoc struct {
	Doc     string
	Options []OptionDoc
}

func (r ResourceDoc) String() string {
	var b bytes.Buffer
	idents := make([]string, 0, len(r.Options))
	for _, opt := range r.Options {
		idents = append(idents, opt.Ident)
	}
	fmt.Fprintf(&b, "Options: %v\n", idents)
	fmt.Fprintf(&b, "%s\n", trunc(r.Doc, 80))
	return b.String()
}

type DomainDoc struct {
	Doc       string
	Resources []ResourceDoc
}

type Definition struct {
	// Name of the definition in the correct format
	// e.g. path:to:mark (without the plus)
	Ident string

	// Target defines on which kind of go expression (Struct etc.) the Definition is appliable
	Targets []target.Target

	// Doc provides documentation for the user to inform about the usage and
	// intention of the definition.
	Doc *OptionDoc

	// DeprecatedInFavorOf points to the marker identifier which should
	// be used instead.
	DeprecatedInFavorOf string

	// output is the type to which parser.Marker has to be converted by a
	// converter.
	Output reflect.Type
}

func (d *Definition) DeprecateInFavorOf(ident string) {
	d.DeprecatedInFavorOf = ident
}

func (d *Definition) IsDeprecated() (string, bool) {
	return d.DeprecatedInFavorOf, d.DeprecatedInFavorOf != ""
}

func (d *Definition) HasDoc() bool {
	return d.Doc != nil
}

func (d *Definition) IsValid() error {
	if err := syntax.Ident(d.Ident); err != nil {
		return err
	}
	if d.Output == nil {
		return fmt.Errorf("output type cannot be nil: %s", d.Ident)
	}
	if len(d.Targets) == 0 {
		return fmt.Errorf("definition has not target defined: %s", d.Ident)
	}
	return nil
}
