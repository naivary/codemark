package testing

import (
	"fmt"
	"math"
	"math/rand/v2"
	"reflect"
	"strconv"
	"unicode"

	"github.com/naivary/codemark/parser"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

const RandLen = -1

// AlmostEqual checks if the two float values a and b are equal with respect to
// some threshold.
func AlmostEqual(a, b float64) bool {
	const float64EqualityThreshold = 1e-5
	return math.Abs(a-b) <= float64EqualityThreshold
}

// NewIdent returns a identifier which can be used for both marker and definition.
// This function should only be used for test purposes. The produced marker is in
// the codemark:testing:* namespace. For custom identifier naming it's recommneded
// to create your own function.
func NewIdent(name string) string {
	return fmt.Sprintf("codemark:testing:%s", name)
}

// RandMarkerWithIdent is the same as RandMarker but allows to set a custom
// identifier for the marker.
func RandMarkerWithIdent(ident string, rtype reflect.Type) (*parser.Marker, error) {
	v := randValue(rtype)
	if v == nil {
		return nil, fmt.Errorf("no value could be generated for given type: %v\n", rtype)
	}
	value := reflect.ValueOf(v)
	m := parser.NewMarker(ident, sdkutil.MarkerKindOf(rtype), value)
	return &m, nil
}

// RandMarker returns a random marker based on the given rtype. The returned
// marker is always valid.
func RandMarker(rtype reflect.Type) (*parser.Marker, error) {
	name := sdkutil.NameFor(rtype)
	ident := NewIdent(name)
	return RandMarkerWithIdent(ident, rtype)
}

// RandName returns a random string which is valid to use for go variables
func RandName() string {
	name := RandString(RandLen)
	for {
		firstLetter := rune(name[0])
		if !unicode.IsDigit(firstLetter) {
			break
		}
		name = RandString(RandLen)
	}
	return name
}

func randValue(rtype reflect.Type) any {
	if !sdkutil.IsSupported(rtype) {
		return nil
	}
	if sdkutil.IsPrimitive(rtype) {
		return randPrimitiveValue(rtype)
	}
	if sdkutil.IsValidSlice(rtype) {
		return randList(rtype.Elem(), RandLen)
	}
	return nil
}

// randPrimitiveValue returns a random value for the given rtype iff rtype is a
// primitive marker type e.g. non LIST.
func randPrimitiveValue(rtype reflect.Type) any {
	if sdkutil.IsInt(rtype) || sdkutil.IsUint(rtype) {
		return randInt(rtype)()
	}
	if sdkutil.IsString(rtype) {
		return RandString(RandLen)
	}
	if sdkutil.IsBool(rtype) {
		return RandBool()
	}
	if sdkutil.IsFloat(rtype) {
		return RandFloat64()
	}
	if sdkutil.IsComplex(rtype) {
		return RandComplex()
	}
	return nil
}

// randList returns a list of len `n` for rtype. For example if rtype is string
// and n is 10 a list of type any will be returned containing 10 random strings.
// If n is <= 0 then the length of the list will be choosen randomly.
func randList(rtype reflect.Type, n int) []any {
	if n <= 0 {
		n = rand.IntN(8) + 1
	}
	// rtype is definetly a supported primitive type which means we can use
	// `randValueFor` to get a correct value.
	values := make([]any, 0, n)
	for range n {
		values = append(values, randPrimitiveValue(rtype))
	}
	return values
}

func randInt(rtype reflect.Type) func() int64 {
	kind := sdkutil.Deref(rtype).Kind()
	maximums := map[reflect.Kind]int64{
		reflect.Int:    math.MaxInt,
		reflect.Int8:   math.MaxInt8,
		reflect.Int16:  math.MaxInt16,
		reflect.Int32:  math.MaxInt32,
		reflect.Int64:  math.MaxInt64,
		reflect.Uint:   math.MaxInt,
		reflect.Uint8:  math.MaxInt8,
		reflect.Uint16: math.MaxInt16,
		reflect.Uint32: math.MaxInt32,
		reflect.Uint64: math.MaxInt64,
	}
	return func() int64 {
		return rand.Int64N(maximums[kind])
	}
}

func RandInt64() int64 {
	typ := reflect.TypeFor[int64]()
	return randInt(typ)()
}

func RandFloat64() float64 {
	const minN = 1
	const maxN = 100
	f := minN + rand.Float64()*(maxN-minN)
	return f
}

func RandBool() bool {
	return RandInt64()%2 == 1
}

func RandString(n int) string {
	if n <= 0 {
		n = 12
	}
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

func RandComplex() complex128 {
	r := rand.IntN(100)
	c := rand.IntN(100)
	compString := fmt.Sprintf("%d+%di", r, c)
	comp, err := strconv.ParseComplex(compString, 128)
	if err != nil {
		panic(err)
	}
	return comp
}
