package registry

import (
	"errors"

	"github.com/naivary/codemark/api/core"
)

var (
	ErrRegistryEmpty = errors.New("the registry contains no definitions. Make sure to add your definitions using the `Define` method")
)

type Registry interface {
	// Define the definition in the registry for future retrieval. It's
	// important to make sure `def.ident` is unique in the Registry.
	Define(def *core.Option) error

	// Get the definition by the unique identiffier
	Get(ident string) (*core.Option, error)

	// DofOf returns the documentation of the definition
	DocOf(ident string) (*core.OptionDoc, error)

	// All returns all Definitions stored in the registry.
	All() map[string]*core.Option
}
