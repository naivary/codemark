package sdk

type Registry interface {
	// Define the definition in the registry for future retrieval. It's
	// important to make sure `def.ident` is unique in the Registry.
	Define(def *Definition) error

	// Get the definition by the unique identiffier
	Get(ident string) (*Definition, error)

	// All returns all Definitions stored in the registry.
	All() map[string]*Definition
}
