package sdk

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/lexer"
)

type DefinitionMaker interface {
	MakeDef(ident string, output reflect.Type, targets ...Target) (*Definition, error)
	MakeDefWithHelp(ident string, output reflect.Type, help *DefinitionHelp, targets ...Target) (*Definition, error)

	MustMakeDef(ident string, output reflect.Type, targets ...Target) *Definition
	MustMakeDefWithHelp(ident string, output reflect.Type, help *DefinitionHelp, targets ...Target) *Definition
}

// TODO: Make it possible to store a Definition in a database
type Definition struct {
	// Name of the definition in the correct format
	// e.g. +path:to:mark
	Ident string

	// Target defines on which type the Definition is appliable
	// e.g. Struct, Package, Field, VAR, CONST etc.
	Targets []Target

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

func (d *Definition) IsValid() error {
	if err := lexer.IsValidIdent(d.Ident); err != nil {
		return err
	}
	if d.Output == nil {
		return fmt.Errorf("output type cannot be nil: %s\n", d.Ident)
	}
	return nil
}
