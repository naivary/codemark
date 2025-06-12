package testing

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/naivary/codemark/parser"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

func RandomMarker(rtype reflect.Type) parser.Marker {
	n := rand.IntN(12) + 1
	typeID := sdkutil.TypeID(rtype)
	ident, found := typeIDToMarkerName[typeID]
	if !found {
		ident = fmt.Sprintf("path:to:%s", typeID)
	}
	v := typeIDToRandValue(typeID, n)
	value := reflect.ValueOf(v)
	return parser.NewMarker(ident, parser.MarkerKindOf(rtype), value)
}

func typeIDToRandValue(typeID string, listLen int) any {
	typeID, _ = strings.CutPrefix(typeID, "ptr.")
	if strings.HasPrefix(typeID, "slice") {
		typeID, _ = strings.CutPrefix(typeID, "slice.")
		return randList(typeID, listLen)
	}
	if strings.HasPrefix(typeID, "string") {
		return randString()
	}
	if strings.HasPrefix(typeID, "bool") {
		return randBool()
	}
	if sdkutil.MatchTypeID(typeID, `int[1-6]{0,2}`) {
		return rand.Int64()
	}
	if sdkutil.MatchTypeID(typeID, `float\d{2}`) {
		return randFloat64()
	}
	if sdkutil.MatchTypeID(typeID, `complex\d{2,3}`) {
		return randComplex()
	}
	return nil
}

func randList(typeID string, listLen int) []any {
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?int\d{0,2}`) {
		return randListT(listLen, randInt64)
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

func randInt64() int64 {
	return rand.Int64()
}

func randFloat64() float64 {
	const minN = 1
	const maxN = 100
	return minN + rand.Float64()*(maxN-minN)
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
