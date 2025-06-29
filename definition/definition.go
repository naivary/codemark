package definition

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/definition/target"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

type Definition struct {
	// Name of the definition in the correct format
	// e.g. +path:to:mark
	Ident string

	// Target defines on which type the Definition is appliable
	// e.g. Struct, Package, Field, VAR, CONST etc.
	Targets []target.Target

	// Doc provides documentation for the user to inform about the usage and
	// intention of the definition.
	Doc string

	// DeprecatedInFavorOf points to the marker identifier which should
	// be used instead.
	DeprecatedInFavorOf *string

	// output is the type to which parser.Marker has to be converted by a
	// converter.
	Output reflect.Type
}

func (d *Definition) DeprecateInFavorOf(marker string) {
	d.DeprecatedInFavorOf = &marker
}

func (d *Definition) IsDeprecated() (*string, bool) {
	return d.DeprecatedInFavorOf, d.DeprecatedInFavorOf != nil
}

func (d *Definition) HasDoc() bool {
	return d.Doc != ""
}

func (d *Definition) IsValid() error {
	if err := sdkutil.IsValidIdent(d.Ident); err != nil {
		return err
	}
	if d.Output == nil {
		return fmt.Errorf("output type cannot be nil: %s\n", d.Ident)
	}
	if len(d.Targets) == 0 {
		return fmt.Errorf("definition has not target defined: %s\n", d.Ident)
	}
	return nil
}
