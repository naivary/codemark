package testing

import (
	"fmt"
	"math"
	"math/rand/v2"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
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

func RandMarkerWithIdent(ident string, rtype reflect.Type) sdk.Marker {
	typeID := sdkutil.TypeIDOf(rtype)
	v := randValueFromTypeID(typeID)
	value := reflect.ValueOf(v)
	return parser.NewMarker(ident, sdk.MarkerKindOf(rtype), value)
}

func RandMarker(rtype reflect.Type) sdk.Marker {
	typeID := sdkutil.TypeIDOf(rtype)
	v := randValueFromTypeID(typeID)
	value := reflect.ValueOf(v)
	return parser.NewMarker(NewIdent(typeID), sdk.MarkerKindOf(rtype), value)
}

func randValueFromTypeID(typeID string) any {
	typeID, _ = strings.CutPrefix(typeID, "ptr.")
	if strings.HasPrefix(typeID, "slice") {
		typeID, _ = strings.CutPrefix(typeID, "slice.")
		return randList(typeID)
	}
	if strings.HasPrefix(typeID, "string") {
		return randString()
	}
	if strings.HasPrefix(typeID, "bool") {
		return randBool()
	}
	if sdkutil.MatchTypeID(typeID, `int[1-6]{0,2}`) {
		return randInt(typeID)()
	}
	if sdkutil.MatchTypeID(typeID, `float\d{2}`) {
		return randFloat64()
	}
	if sdkutil.MatchTypeID(typeID, `complex\d{2,3}`) {
		return randComplex()
	}
	return nil
}

func randList(typeID string) []any {
	listLen := rand.IntN(8) + 1
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?int\d{0,2}`) {
		return randListT(listLen, randInt(typeID))
	}
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?float\d{2}`) {
		return randListT(listLen, randFloat64)
	}
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?complex\d{2,3}`) {
		return randListT(listLen, randComplex)
	}
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?bool`) {
		return randListT(listLen, randBool)
	}
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?string`) {
		return randListT(listLen, randString)
	}
	return nil
}

func randInt(typeID string) func() int64 {
	typeID, _ = strings.CutPrefix(typeID, "ptr.")
	maxs := map[sdk.TypeID]int64{
		"int":    math.MaxInt,
		"int8":   math.MaxInt8,
		"int16":  math.MaxInt16,
		"int32":  math.MaxInt32,
		"int64":  math.MaxInt64,
		"uint":   math.MaxInt,
		"uint8":  math.MaxInt8,
		"uint16": math.MaxInt16,
		"uint32": math.MaxInt32,
		"uint64": math.MaxInt64,
	}
	return func() int64 {
		return rand.Int64N(maxs[typeID])
	}
}

func randInt64() int64 {
	return randInt("int64")()
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

func randListT[T any](n int, rnd func() T) []any {
	vals := make([]any, 0, n)
	for range n {
		val := rnd()
		vals = append(vals, val)
	}
	return vals
}
