package registry

import (
	"errors"

	"github.com/naivary/codemark/definition"
)

var (
	ErrRegistryEmpty = errors.New("the registry contains no definitions. Make sure to add your definitions using the `Define` method")
)

type Registry interface {
	// Define the definition in the registry for future retrieval. It's
	// important to make sure `def.ident` is unique in the Registry.
	Define(def *definition.Definition) error

	// Get the definition by the unique identiffier
	Get(ident string) (*definition.Definition, error)

	// DofOf returns the documentation of the definition
	DocOf(ident string) (string, error)

	// All returns all Definitions stored in the registry.
	All() map[string]*definition.Definition
}
