package registry

import (
	"errors"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
)

var ErrRegistryEmpty = errors.New(
	"the registry contains no definitions. Make sure to add your definitions using the `Define` method",
)

// TODO: move this under api/
type Registry interface {
	// Define defines the option in the registry. Options must be unique.
	Define(opt *optionv1.Option) error

	// Get the definition by the unique identiffier
	Get(ident string) (*optionv1.Option, error)

	// DofOf returns the documentation of the definition
	DocOf(ident string) (*docv1.Option, error)

	// All returns all Definitions stored in the registry.
	All() map[string]*optionv1.Option
}
