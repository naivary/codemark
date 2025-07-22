package registry

import (
	"errors"

	"github.com/naivary/codemark/api/core"
)

var ErrRegistryEmpty = errors.New(
	"the registry contains no definitions. Make sure to add your definitions using the `Define` method",
)

type Registry interface {
	// Define defines the option in the registry. Options must be unique.
	Define(opt *core.Option) error

	// Get the definition by the unique identiffier
	Get(ident string) (*core.Option, error)

	// DofOf returns the documentation of the definition
	DocOf(ident string) (*core.OptionDoc, error)

	// All returns all Definitions stored in the registry.
	All() map[string]*core.Option
}
