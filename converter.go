package codemark

import "github.com/naivary/codemark/converter"

// NewConvManager returns a new converter manager. One case you might need a
// converter manager is for adding your custom converters in the load process.
var NewConvManager = converter.NewManager
