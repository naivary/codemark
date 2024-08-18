package main

import "reflect"

func MakeDef(name string, t Target, output reflect.Type) *Definition {
	def := &Definition{
		name:   name,
		target: t,
		output: output,
	}
	def.resolveKind()
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
	name string

	// The output type to which the value
	// of the marker will be mapped to
	output reflect.Type

	// TargetType defines on which type of
	// target it can be applied e.g. constants,
	// functions, types, variables etc.
	target Target

	help *DefinitionHelp

	deprecatedInFavorOf *string

	kind reflect.Kind

	isPointer bool

	sliceKind reflect.Kind
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

func (d *Definition) Name() string {
	return d.name
}

func (d *Definition) Type() reflect.Type {
	return d.output
}

func (d *Definition) Target() Target {
	return d.target
}

func (d *Definition) Help() *DefinitionHelp {
	return d.help
}

func (d *Definition) resolveKind() {
	d.kind, d.isPointer = resolvePtr(d)
	if d.kind == reflect.Slice {
		d.sliceKind = d.output.Elem().Kind()
	}
}
