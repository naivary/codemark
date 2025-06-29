package sdk

import (
	"reflect"

	"github.com/naivary/codemark/parser"
)

type Converter interface {
	// Name is a human friendly representation of the converter. It may or may
	// not be unique depending on the Converter. It is recommended to prefix the
	// name of the converter with your project. For example the builtin
	// converters are prefixed with `codemark.*`.
	Name() string

	// SupportedTypes returns the supported types for which the converter can be
	// used.
	SupportedTypes() []reflect.Type

	// CanConvert validates if the conversion of the marker to `def.output` is
	// possible. You can be sure that the convert is choosen correctly by the
	// ConverterManager and do not have to check if the `def.Output` is
	// convertible using this converter.
	CanConvert(m parser.Marker, def *Definition) error

	// Convert converts the marker to `def.Output`
	Convert(m parser.Marker, def *Definition) (reflect.Value, error)
}

// ConverterManager is responsible for managing the workflow of converting a
// marker to a definition.
type ConverterManager interface {
	// GetConverter returns the converter for the given reflect.Type. If none is
	// found an error will be returned.
	GetConverter(from reflect.Type) (Converter, error)

	// AddConverter adds the converter to the manager.
	AddConverter(conv Converter) error

	// Convert converts the marker with respect to the target to identified
	// definition in the registry.
	Convert(mrk parser.Marker, target Target) (any, error)

	// ParseDefs returns all definitions found in the `doc` with respect to the
	// target.
	ParseDefs(doc string, t Target) (map[string][]any, error)
}
