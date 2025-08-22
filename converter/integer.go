package converter

import (
	"fmt"
	"reflect"
	"slices"
	"time"

	convv1 "github.com/naivary/codemark/api/converter/v1"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/rtypeutil"
)

var _ convv1.Converter = (*intConverter)(nil)

type intConverter struct {
	name string
}

func NewInteger() convv1.Converter {
	return &intConverter{
		name: "integer",
	}
}

func (i *intConverter) Name() string {
	return NewName(_codemark, i.name)
}

func (i *intConverter) SupportedTypes() []reflect.Type {
	types := []any{
		int(0),
		int8(0),
		int16(0),
		// rune
		int32(0),
		time.Duration(0),
		int64(0),
		uint(0),
		// byte
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),

		// pointer
		new(int),
		new(int8),
		new(int16),
		// *rune
		new(int32),
		new(int64),
		new(time.Duration),
		new(uint),
		// *byte
		new(uint8),
		new(uint16),
		new(uint32),
		new(uint64),
	}
	supported := make([]reflect.Type, 0, len(types))
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		supported = append(supported, rtype)
	}
	return supported
}

func (i *intConverter) CanConvert(m marker.Marker, to reflect.Type) error {
	mkind := m.Kind
	out := rtypeutil.Deref(to)
	if mkind == marker.INT {
		return nil
	}
	if mkind == marker.STRING &&
		slices.Contains([]reflect.Kind{_byte, _rune, _duration}, out.Kind()) {
		return nil
	}
	return fmt.Errorf(
		"marker kind of `%s` cannot be converted to a int. valid options are: %s;%s",
		mkind,
		marker.INT,
		marker.STRING,
	)
}

func (i *intConverter) Convert(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	mkind := m.Kind
	if i.isInteger(to, mkind) {
		return i.integer(m, to)
	}
	if i.isUint(to, mkind) {
		return i.uinteger(m, to)
	}
	if i.isDuration(to, mkind) {
		return i.duration(m, to)
	}
	if i.isByte(to, mkind) {
		return i.bytee(m, to)
	}
	if i.isRune(to, mkind) {
		return i.runee(m, to)
	}
	return _rvzero, fmt.Errorf("cannot convert %s to %v", m.Ident, to)
}

func (i *intConverter) integer(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	n := m.Value.Int()
	if i.isOverflowingInt(to, n) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`", m.String(), to)
	}
	return ConvertTo(m.Value, to)
}

func (i *intConverter) uinteger(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	n := m.Value.Int()
	if i.isOverflowingUint(to, uint64(n)) {
		return _rvzero, fmt.Errorf("overflow converting `%s` to `%v`", m.String(), to)
	}
	return ConvertTo(m.Value, to)
}

func (i *intConverter) runee(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	v := m.Value.String()
	if len(v) > 1 {
		return _rvzero, fmt.Errorf(
			"marker value cannot be bigger than 2 chars for rune conversion: %s",
			v,
		)
	}
	rvalue := reflect.ValueOf(rune(v[0]))
	return ConvertTo(rvalue, to)
}

func (i *intConverter) bytee(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	v := m.Value.String()
	if len(v) > 1 {
		return _rvzero, fmt.Errorf("value of marker is bigger than 2: %s", v)
	}
	bvalue := reflect.ValueOf(byte(v[0]))
	return ConvertTo(bvalue, to)
}

func (i *intConverter) duration(m marker.Marker, to reflect.Type) (reflect.Value, error) {
	duration, err := time.ParseDuration(m.Value.String())
	if err != nil {
		return _rvzero, err
	}
	v := reflect.ValueOf(duration)
	return ConvertTo(v, to)
}

func (i *intConverter) isOverflowingInt(out reflect.Type, n int64) bool {
	return rtypeutil.Deref(out).OverflowInt(n)
}

func (i *intConverter) isOverflowingUint(out reflect.Type, n uint64) bool {
	return rtypeutil.Deref(out).OverflowUint(n)
}

func (i *intConverter) isInteger(rtype reflect.Type, mkind marker.Kind) bool {
	return rtypeutil.IsInt(rtype) && mkind == marker.INT
}

func (i *intConverter) isUint(rtype reflect.Type, mkind marker.Kind) bool {
	return rtypeutil.IsUint(rtype) && mkind == marker.INT
}

func (i *intConverter) isByte(rtype reflect.Type, mkind marker.Kind) bool {
	rtype = rtypeutil.Deref(rtype)
	return rtype.Kind() == _byte && mkind == marker.STRING
}

func (i *intConverter) isRune(rtype reflect.Type, mkind marker.Kind) bool {
	rtype = rtypeutil.Deref(rtype)
	return rtype.Kind() == _rune && mkind == marker.STRING
}

func (i *intConverter) isDuration(rtype reflect.Type, mkind marker.Kind) bool {
	rtype = rtypeutil.Deref(rtype)
	return rtype.Kind() == _duration && mkind == marker.STRING
}
