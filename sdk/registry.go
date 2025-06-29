package sdk

import "github.com/naivary/codemark/definition"

type Registry interface {
	// Define the definition in the registry for future retrieval. It's
	// important to make sure `def.ident` is unique in the Registry.
	Define(def *definition.Definition) error

	// Get the definition by the unique identiffier
	Get(ident string) (*definition.Definition, error)

	// All returns all Definitions stored in the registry.
	All() map[string]*definition.Definition
}
