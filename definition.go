package main

import "reflect"

type DefinitionHelp struct {
	Category string

	Summary string

	DeprecatedInFavorOf bool
}

type Definition struct {
	// Name of the definition in the correct format
	// e.g. path:to:mark
	Name string

	Output reflect.Type

	TargetType Target

	Help *DefinitionHelp
}

func MakeDefinition(name string, targetType Target, output any) *Definition {
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
