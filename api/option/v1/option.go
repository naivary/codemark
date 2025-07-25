package v1

import (
	"bytes"
	"fmt"
	"reflect"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	"github.com/naivary/codemark/internal/console"
)

type Option struct {
	// Ident is the identifier of the option e.g. domain:resource:option
	Ident string

	// Target defines on which kind of go expression (Struct etc.) the option is appliable
	Targets []Target

	// Doc provides documentation for the user to inform about the usage and
	// intention of the option.
	Doc *docv1.Option

	// DeprecatedInFavorOf points to the marker identifier which should
	// be used instead.
	DeprecatedInFavorOf string

	// Output type to convert the marker to.
	Output reflect.Type

	// Whether this option is unique
	IsUnique bool
}

func (o *Option) DeprecateInFavorOf(ident string) {
	o.DeprecatedInFavorOf = ident
}

func (o *Option) IsDeprecated() bool {
	return o.DeprecatedInFavorOf != ""
}

func (o *Option) HasDoc() bool {
	return o.Doc != nil
}

func (o *Option) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "Ident: %s\n", o.Ident)
	fmt.Fprintf(&b, "Default: %s\n", o.Doc.Default)
	fmt.Fprintf(&b, "%s", console.Trunc(o.Doc.Desc, 80))
	return b.String()
}
