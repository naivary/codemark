package sdk

type Registry interface {
	// Define the definition in the registry for future retrieval. It's
	// important to make sure `def.idn` is unique in the Registry.
	Define(def *Definition) error

	// Get the definition by the unique identiffier
	Get(idn string) (*Definition, error)

	All() map[string]*Definition
}
