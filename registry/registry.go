package registry

import (
	"errors"

	optionapi "github.com/naivary/codemark/api/option"
)

var ErrRegistryEmpty = errors.New(
	"the registry contains no definitions. Make sure to add your definitions using the `Define` method",
)

type Registry interface {
	// Define defines the option in the registry. Options must be unique.
	Define(opt *optionapi.Option) error

	// Get the definition by the unique identiffier
	Get(ident string) (*optionapi.Option, error)

	// DofOf returns the documentation of the definition
	DocOf(ident string) (*optionapi.OptionDoc, error)

	// All returns all Definitions stored in the registry.
	All() map[string]*optionapi.Option

	Merge(reg Registry) error
}
