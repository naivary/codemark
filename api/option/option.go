package option

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/naivary/codemark/internal/console"
)

type Option struct {
	// Ident is the identifier of the option e.g. domain:resource:option
	Ident string

	// Target defines on which kind of go expression (Struct etc.) the option is appliable
	Targets []Target

	// Doc provides documentation for the user to inform about the usage and
	// intention of the option.
	Doc *OptionDoc

	// DeprecatedInFavorOf points to the marker identifier which should
	// be used instead.
	DeprecatedInFavorOf string

	// Output type to convert the marker to.
	Output reflect.Type
}

func (o *Option) DeprecateInFavorOf(ident string) {
	o.DeprecatedInFavorOf = ident
}

func (o *Option) IsDeprecated() (string, bool) {
	return o.DeprecatedInFavorOf, o.DeprecatedInFavorOf != ""
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
