package codemark

import (
	"reflect"
)

func MakeDef(idn string, t Target, output any) *Definition {
	outputType := reflect.TypeOf(output)
	def := &Definition{
		Ident:      idn,
		Target:     t,
		output:     outputType,
		underlying: underlying(outputType),
	}
	def.kind = def.nonPtrKind()
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

	underlying reflect.Type

	kind reflect.Kind

	// The output type to which the value
	// of the marker will be converted to
	output reflect.Type
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

func (d *Definition) typ() reflect.Type {
	if d.sliceType() != nil {
		return d.sliceType()
	}
	return d.underlying
}

func (d *Definition) nonPtrType() reflect.Type {
	typ := d.typ()
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ
}

func (d *Definition) nonPtrKind() reflect.Kind {
	kind := d.output.Kind()
	if kind != reflect.Pointer {
		return kind
	}
	return d.output.Elem().Kind()
}

func (d *Definition) sliceType() reflect.Type {
	if d.kind != reflect.Slice {
		return nil
	}
	typ := d.output.Elem()
	if typ.Kind() == reflect.Slice {
		return typ.Elem()
	}
	return typ
}

func (d *Definition) sliceKind() reflect.Kind {
	typ := d.sliceType()
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ.Kind()
}
