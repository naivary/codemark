package v1

import (
	"reflect"

	docv1 "github.com/naivary/codemark/api/doc/v1"
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

	// Type type to convert the marker to.
	Type reflect.Type

	// Whether this option is unique
	IsUnique bool

	// The Default value of the option. If nil the option is taken as required
	Default any
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

func (o *Option) IsRequired() bool {
	return o.Default == nil
}
