package converter

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	coreapi "github.com/naivary/codemark/api/core"
	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/typeutil"
)

const (
	_rune     = reflect.Int32
	_byte     = reflect.Uint8
	_duration = reflect.Int64
)

const _codemark = "codemark"

// _rvzero is the zero value for a reflect.Value used for convenience
var _rvzero = reflect.Value{}

func Get(rtype reflect.Type) converter.Converter {
	if typeutil.IsBool(rtype) {
		return Bool()
	}
	if typeutil.IsString(rtype) {
		return String()
	}
	if typeutil.IsInt(rtype) || typeutil.IsUint(rtype) {
		return Integer()
	}
	if typeutil.IsFloat(rtype) {
		return Float()
	}
	if typeutil.IsComplex(rtype) {
		return Complex()
	}
	if typeutil.IsAny(rtype) {
		return Any()
	}
	return nil
}

func isCorrectTarget(opt coreapi.Option, t coreapi.Target) bool {
	return !slices.Contains(opt.Targets, t) && !slices.Contains(opt.Targets, coreapi.TargetAny)
}

// isValidName checks if the choosen name of a custom converter is following the
// convention of prefixing the name with the project name and that the project
// name is not "codemark".
func isValidName(name string) error {
	if strings.HasPrefix(name, _codemark) {
		return fmt.Errorf(`the name of your custom converter cannot start with "codemark" because it is reserved for the builtin converters: %s`, name)
	}
	if len(strings.Split(name, ".")) != 2 {
		return fmt.Errorf(`the name of your custom converter has to be seperated with "%s" and must be composed of two segments e.g. "codemark.integer"`, ".")
	}
	return nil
}

func isTypeT[T any](to reflect.Type) bool {
	var from T
	to = typeutil.Deref(to)
	return to.ConvertibleTo(reflect.TypeOf(from))
}
