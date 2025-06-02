package codemark

import (
	"reflect"
)

func MakeDef(idn string, t Target, output any) *Definition {
	def := &Definition{
		Ident:  idn,
		Target: t,
		output: reflect.TypeOf(output),
	}
	return def
}

func MakeDefWithHelp(name string, t Target, output reflect.Type, help *DefinitionHelp) *Definition {
	def := MakeDef(name, t, output)
	def.Help = help
	return def
}

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

	output reflect.Type
}

type DefinitionHelp struct {
	Category string

	Description string
}

func (d *Definition) Output() reflect.Type {
	return d.output
}

func (d *Definition) DeprecateInFavorOf(marker string) {
	d.DeprecatedInFavorOf = &marker
}

func (d *Definition) IsDeprecated() (*string, bool) {
	return d.DeprecatedInFavorOf, d.DeprecatedInFavorOf != nil
}
