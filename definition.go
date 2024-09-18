package codemark

import (
	"reflect"
)

func MakeDef(idn string, t Target, output reflect.Type) *Definition {
	def := &Definition{
		identifier: idn,
		target:     t,
		output:     output,
		underlying: underlying(output),
	}
	def.kind = def.nonPtrKind()
	return def
}

func MakeDefWithHelp(name string, t Target, output reflect.Type, help *DefinitionHelp) *Definition {
	def := MakeDef(name, t, output)
	def.help = help
	return def
}

type Definition struct {
	// Name of the definition in the correct format
	// e.g. +path:to:mark
	identifier string

	// The output type to which the value
	// of the marker will be converted to
	output reflect.Type

	// TargetType defines on which type of
	// target it can be applied e.g. constants,
	// functions, types, variables etc.
	target Target

	help *DefinitionHelp

	deprecatedInFavorOf *string

	underlying reflect.Type

	kind reflect.Kind
}

type DefinitionHelp struct {
	Category string

	Description string
}

func (d *Definition) DeprecateInFavorOf(marker string) {
	d.deprecatedInFavorOf = &marker
}

func (d *Definition) IsDeprecated() (*string, bool) {
	return d.deprecatedInFavorOf, d.deprecatedInFavorOf == nil
}

func (d *Definition) Idn() string {
	return d.identifier
}

func (d *Definition) Output() reflect.Type {
	return d.output
}

func (d *Definition) Target() Target {
	return d.target
}

func (d *Definition) Help() *DefinitionHelp {
	return d.help
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
