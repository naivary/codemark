package main

import "reflect"

type DefinitionHelp struct {
	Category string

	Summary string
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

    // Help provides help information about 
    // the marker
	Help *DefinitionHelp

	DeprecatedInFavorOf *string
}

func MakeDef(name string, targetType Target, output any) *Definition{
	return &Definition{
		Name:       name,
		TargetType: targetType,
		Output:     reflect.TypeOf(output),
	}
}

func MakeDefWithHelp(name string, targetType Target, output any, help *DefinitionHelp) *Definition {
	return &Definition{
		Name:       name,
		TargetType: targetType,
		Output:     reflect.TypeOf(output),
		Help:       help,
	}
}
