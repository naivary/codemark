package converter

import (
	"reflect"

	"github.com/naivary/codemark/converter"
	"github.com/naivary/codemark/typeutil"
)

const (
	_rune = reflect.Int32
	_byte = reflect.Uint8
)

const _codemark = "codemark"

var (
	// _rvzero is the zero value for a reflect.Value used for convenience
	_rvzero = reflect.Value{}
)

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
