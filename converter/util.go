package converter

import (
	"fmt"
	"reflect"

	"github.com/naivary/codemark/typeutil"
)

const NameSep = "."

// NewName is returning a valid converter name. The convention is to prefix
// every converter with your project name, followed by a custom name for the
// converter seperated by a dot.
func NewName(proj, conv string) string {
	return fmt.Sprintf("%s.%s", proj, conv)
}

// ConvertTo converts the value `v` to the Type `typ` handling pointer
// dereferencing and other inconveniences.
func ConvertTo(v reflect.Value, to reflect.Type) (reflect.Value, error) {
	isPtr := typeutil.IsPointer(to)
	outputType := to
	// need to dereference type to create the correct variable using
	// `reflect.New`. Otherwise .Set wont work.
	if isPtr {
		outputType = outputType.Elem()
	}
	out := reflect.New(outputType)
	out.Elem().Set(v.Convert(outputType))
	if isPtr {
		return out.Convert(to), nil
	}
	return out.Elem(), nil
}
