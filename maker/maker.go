package maker

import (
	"reflect"

	"github.com/naivary/codemark/sdk"
)

var _ sdk.DefinitionMaker = (*defMaker)(nil)

type defMaker struct{}

func New() sdk.DefinitionMaker {
	return &defMaker{}
}

func (d defMaker) MakeDef(ident string, output reflect.Type, targets ...sdk.Target) (*sdk.Definition, error) {
	return MakeDef(ident, output, targets...)
}

func (d defMaker) MakeDefWithHelp(ident string, output reflect.Type, help *sdk.DefinitionHelp, targets ...sdk.Target) (*sdk.Definition, error) {
	return MakeDefWithHelp(ident, output, help, targets...)
}

func (d defMaker) MustMakeDef(ident string, output reflect.Type, targets ...sdk.Target) *sdk.Definition {
	return MustMakeDef(ident, output, targets...)
}

func (d defMaker) MustMakeDefWithHelp(ident string, output reflect.Type, help *sdk.DefinitionHelp, targets ...sdk.Target) *sdk.Definition {
	return MustMakeDefWithHelp(ident, output, help, targets...)
}
