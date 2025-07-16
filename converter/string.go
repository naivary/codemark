package converter

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/parser/marker"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var _ sdk.Converter = (*stringConverter)(nil)

type stringConverter struct {
	name string
}

func String() sdk.Converter {
	return &stringConverter{
		name: "string",
	}
}

func (s *stringConverter) Name() string {
	return sdkutil.NewConvName(_codemark, s.name)
}

func (s *stringConverter) SupportedTypes() []reflect.Type {
	types := []any{
		string(""),
		// pointer
		new(string),
	}
	supported := make([]reflect.Type, 0, len(types))
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		supported = append(supported, rtype)
	}
	return supported
}

func (s *stringConverter) CanConvert(m marker.Marker, to reflect.Type) error {
	if m.Kind != marker.STRING {
		return fmt.Errorf("marker kind of `%s` cannot be converted to a string. valid option is: %s\n", m.Kind, marker.STRING)
	}
	return nil
}

func (s *stringConverter) Convert(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	return sdkutil.ConvertTo(m.Value, to)
}
