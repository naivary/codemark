package sdk

import (
	"reflect"
)

// TODO: Make it possible to store a Definition in a database
type Definition struct {
	// Name of the definition in the correct format
	// e.g. +path:to:mark
	Ident string

	// Target defines on which type the Definition is appliable
	// e.g. Struct, Package, Field, VAR, CONST etc.
	Target Target

	// Help provides user-defined documentation for the definition
	Help *DefinitionHelp

	// DeprecatedInFavorOf points to the marker identifier which should
	// be used instead.
	DeprecatedInFavorOf *string

	// output is the type to which parser.Marker has to be converted by a
	// converter.
	Output reflect.Type
}

type DefinitionHelp struct {
	Category string

	Description string
}

func (d *Definition) DeprecateInFavorOf(marker string) {
	d.DeprecatedInFavorOf = &marker
}

func (d *Definition) IsDeprecated() (*string, bool) {
	return d.DeprecatedInFavorOf, d.DeprecatedInFavorOf != nil
}
