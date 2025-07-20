package markertest

import (
	"fmt"
	randv2 "math/rand/v2"
	"reflect"

	"github.com/naivary/codemark/internal/rand"
	"github.com/naivary/codemark/marker"
	"github.com/naivary/codemark/typeutil"
)

// NewIdent returns a identifier which can be used for both marker and definition.
// This function should only be used for test purposes. The produced marker is in
// the codemark:testing:* namespace. For custom identifier naming it's recommneded
// to create your own function.
func NewIdent(name string) string {
	return fmt.Sprintf("codemark:testing:%s", name)
}

// RandMarkerWithIdent is the same as RandMarker but allows to set a custom
// identifier for the marker.
func RandMarkerWithIdent(ident string, rtype reflect.Type) (*marker.Marker, error) {
	v := randValue(rtype)
	if v == nil {
		return nil, fmt.Errorf("no value could be generated for given type: %v\n", rtype)
	}
	value := reflect.ValueOf(v)
	m := marker.New(ident, marker.KindOf(rtype), value)
	return &m, nil
}

// RandMarker returns a random marker based on the given rtype. The returned
// marker is always valid if not error is returned.
func RandMarker(rtype reflect.Type) (*marker.Marker, error) {
	name := typeutil.NameFor(rtype)
	ident := NewIdent(name)
	return RandMarkerWithIdent(ident, rtype)
}

// randValue returns a valid marker value for the given rtype.
func randValue(rtype reflect.Type) any {
	if !typeutil.IsSupported(rtype) {
		return nil
	}
	if typeutil.IsPrimitive(rtype) {
		return randPrimitiveValue(rtype)
	}
	if typeutil.IsValidSlice(rtype) {
		return randList(rtype.Elem(), rand.RandLen)
	}
	return nil
}

// randPrimitiveValue returns a random value for the given rtype iff rtype is a
// primitive marker type e.g. non LIST.
func randPrimitiveValue(rtype reflect.Type) any {
	if typeutil.IsInt(rtype) || typeutil.IsUint(rtype) {
		return rand.IntFromType(rtype)()
	}
	if typeutil.IsString(rtype) {
		return rand.String(rand.RandLen)
	}
	if typeutil.IsBool(rtype) {
		return rand.Bool()
	}
	if typeutil.IsFloat(rtype) {
		return rand.Float64()
	}
	if typeutil.IsComplex(rtype) {
		return rand.Complex()
	}
	if typeutil.IsAny(rtype) {
		return randPrimitiveValue(randPrimitiveType())
	}
	return nil
}

func randPrimitiveType() reflect.Type {
	i := randv2.IntN(5)
	switch i {
	case 0:
		return reflect.TypeFor[int]()
	case 1:
		return reflect.TypeFor[string]()
	case 2:
		return reflect.TypeFor[bool]()
	case 3:
		return reflect.TypeFor[float64]()
	case 4:
		return reflect.TypeFor[complex128]()
	}
	return nil
}

// randList returns a list of length `n` for rtype. For example if rtype is string
// and n is 10 a list of type any will be returned containing 10 random strings.
// If n is <= 0 then the length of the list will be choosen randomly.
func randList(rtype reflect.Type, n int) []any {
	if n <= 0 {
		n = randv2.IntN(8) + 1
	}
	values := make([]any, 0, n)
	for range n {
		values = append(values, randPrimitiveValue(rtype))
	}
	return values
}
