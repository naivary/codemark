package codemark

import (
	"fmt"
	"reflect"
)

const (
	_rune       = reflect.Int32
	_byte       = reflect.Uint8
	_namePrefix = "codemark"
)

var (
	_rvzero = reflect.Value{}
)

func buildName(name string) string {
	return fmt.Sprintf("%s.%s", _namePrefix, name)
}
