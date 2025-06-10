package sdk

import "errors"

var (
	ErrRegistryEmpty = errors.New("the registry contains no definitions. Make sure to add your definitions using the `Define` method")
)
