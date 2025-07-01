package converter

import (
	"reflect"
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
