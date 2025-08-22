package converter

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	optionv1 "github.com/naivary/codemark/api/option/v1"
	"github.com/naivary/codemark/rtypeutil"
)

const (
	_rune     = reflect.Int32
	_byte     = reflect.Uint8
	_duration = reflect.Int64
)

const _codemark = "codemark"

// _rvzero is the zero value for a reflect.Value used for convenience
var _rvzero = reflect.Value{}

// NewName is returning a valid converter name. The convention is to prefix
// every converter with your project name, followed by a custom name for the
// converter seperated by a dot.
func NewName(proj, conv string) string {
	return fmt.Sprintf("%s.%s", proj, conv)
}

// ConvertTo converts the value `v` to the Type `typ` handling pointer
// dereferencing and other inconveniences.
func ConvertTo(v reflect.Value, to reflect.Type) (reflect.Value, error) {
	isPtr := rtypeutil.IsPointer(to)
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

func isCorrectTarget(opt optionv1.Option, t optionv1.Target) bool {
	return !slices.Contains(opt.Targets, t) && !slices.Contains(opt.Targets, optionv1.TargetAny)
}

// isValidName checks if the choosen name of a custom converter is following the
// convention of prefixing the name with the project name and that the project
// name is not "codemark".
func isValidName(name string) error {
	if strings.HasPrefix(name, _codemark) {
		return fmt.Errorf(
			`the name of your custom converter cannot start with "codemark" because it is reserved for the builtin converters: %s`,
			name,
		)
	}
	if len(strings.Split(name, ".")) != 2 {
		return fmt.Errorf(
			`the name of your custom converter has to be seperated with "%s" and must be composed of two segments e.g. "codemark.integer"`,
			".",
		)
	}
	return nil
}

func isTypeT[T any](to reflect.Type) bool {
	var from T
	to = rtypeutil.Deref(to)
	return to.ConvertibleTo(reflect.TypeOf(from))
}
