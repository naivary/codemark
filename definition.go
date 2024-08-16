package main

import "reflect"

type DefinitionHelp struct {
	Category string

	Description string
}

type Definition struct {
	// Name of the definition in the correct format
	// e.g. +path:to:mark
	Name string

	// The output type to which the value
	// of the marker will be mapped to
	Output reflect.Type

	// TargetType defines on which type of
	// target it can be applied e.g. constants,
	// functions, types, variables etc.
	TargetType Target

	Help *DefinitionHelp

	DeprecatedInFavorOf *string
}

func (d *Definition) Deprecate(marker string) {
	d.DeprecatedInFavorOf = &marker
}

func (d *Definition) IsDeprecated() bool {
	return d.DeprecatedInFavorOf != nil
}

func MakeDef(name string, targetType Target, output reflect.Type) *Definition {
	return &Definition{
		Name:       name,
		TargetType: targetType,
		Output:     output,
	}
}

func MakeDefWithHelp(name string, targetType Target, output reflect.Type, help *DefinitionHelp) *Definition {
	return &Definition{
		Name:       name,
		TargetType: targetType,
		Output:     output,
		Help:       help,
	}
}
