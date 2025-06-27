package testing

import (
	"fmt"
	"math"
	"math/rand/v2"
	"reflect"
	"strconv"
	"time"

	"github.com/naivary/codemark/parser"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

// AlmostEqual checks if the two float values a and b are equal with respect to
// some threshold.
func AlmostEqual(a, b float64) bool {
	const float64EqualityThreshold = 1e-5
	return math.Abs(a-b) <= float64EqualityThreshold
}

// NewIdent returns a identifier which can be used for both marker and
// definition. This function should only be used for test purposes and if the
// marker should be included in the codemark:testing:* namespace. For custom
// identifier naming it's recommneded to create your own function.
func NewIdent(typeID string) string {
	return fmt.Sprintf("codemark:testing:%s", typeID)
}

func RandMarkerWithIdent(ident string, rtype reflect.Type) *parser.Marker {
	v := randValueFor(rtype)
	value := reflect.ValueOf(v)
	m := parser.NewMarker(ident, sdkutil.MarkerKindOf(rtype), value)
	return &m
}

func RandMarker(rtype reflect.Type) *parser.Marker {
	v := randValueFor(rtype)
	value := reflect.ValueOf(v)
	m := parser.NewMarker(NewIdent(typeID), sdkutil.MarkerKindOf(rtype), value)
	return &m
}

func randValueFor(rtype reflect.Type) any {
	if sdkutil.IsInt(rtype) {
		return randInt(rtype)
	}
	if sdkutil.IsString(rtype) {
		return randString()
	}
	if sdkutil.IsBool(rtype) {
		return randBool()
	}
	if sdkutil.IsFloat(rtype) {
		return randFloat64()
	}
	if sdkutil.IsComplex(rtype) {
		return randComplex()
	}
	if sdkutil.IsValidSlice(rtype) {
		return randList(rtype.Elem())
	}
	return nil
}

func randList(rtype reflect.Type) []any {
	n := rand.IntN(8) + 1
	// rtype is definetly a supported primitive type which means we can use
	// `randValueFor` to get a correct value.
	values := make([]any, 0, n)
	for range n {
		values = append(values, randValueFor(rtype))
	}
	return values
}

func randInt(rtype reflect.Type) func() int64 {
	kind := sdkutil.Deref(rtype).Kind()
	maxs := map[reflect.Kind]int64{
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
		return rand.Int64N(maxs[kind])
	}
}

func randInt64() int64 {
	typ := reflect.TypeFor[int64]()
	return randInt(typ)()
}

func randFloat64() float64 {
	const minN = 1
	const maxN = 100
	f := minN + rand.Float64()*(maxN-minN)
	return f
}

func randBool() bool {
	return rand.Int64N(time.Now().Unix())%2 == 1
}

func randString() string {
	const n = 12
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

func randComplex() complex128 {
	r := rand.IntN(100)
	c := rand.IntN(100)
	compString := fmt.Sprintf("%d+%di", r, c)
	comp, err := strconv.ParseComplex(compString, 128)
	if err != nil {
		panic(err)
	}
	return comp
}
