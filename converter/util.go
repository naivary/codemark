package converter

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	coreapi "github.com/naivary/codemark/api/core"
	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

const (
	_rune     = reflect.Int32
	_byte     = reflect.Uint8
	_codemark = "codemark"
)

var (
	// _rvzero is the zero value for a reflect.Value used for convenience
	_rvzero = reflect.Value{}
)

func Get(rtype reflect.Type) sdk.Converter {
	if sdkutil.IsBool(rtype) {
		return Bool()
	}
	if sdkutil.IsString(rtype) {
		return String()
	}
	if sdkutil.IsInt(rtype) || sdkutil.IsUint(rtype) {
		return Integer()
	}
	if sdkutil.IsFloat(rtype) {
		return Float()
	}
	if sdkutil.IsComplex(rtype) {
		return Complex()
	}
	if sdkutil.IsAny(rtype) {
		return Any()
	}
	return nil
}

func isCorrectTarget(opt coreapi.Option, t coreapi.Target) bool {
	return !(slices.Contains(opt.Targets, t) || slices.Contains(opt.Targets, coreapi.TargetAny))
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
