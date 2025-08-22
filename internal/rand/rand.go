package rand

import (
	"fmt"
	"math"
	"math/rand/v2"
	"reflect"
	"strconv"
	"unicode"

	"github.com/naivary/codemark/rtypeutil"
)

const RandLen = -1

func Int64() int64 {
	typ := reflect.TypeFor[int64]()
	return IntFromType(typ)()
}

// IntFromType returns a random integer respecting the MaxInt of the type e.g.
// if a type of Int8 is passed the random integer cannot be greater than
// math.MaxInt8. This allows to assure no overflow will occur.
func IntFromType(rtype reflect.Type) func() int64 {
	kind := rtypeutil.Deref(rtype).Kind()
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

func Float64() float64 {
	const minN = 1
	const maxN = 100
	f := minN + rand.Float64()*(maxN-minN)
	return f
}

func Bool() bool {
	return Int64()%2 == 1
}

func String(n int) string {
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

func Complex() complex128 {
	r := rand.IntN(100)
	c := rand.IntN(100)
	compString := fmt.Sprintf("%d+%di", r, c)
	comp, err := strconv.ParseComplex(compString, 128)
	if err != nil {
		panic(err)
	}
	return comp
}

// GoIdent returns a random string which can be used as a go identifier.
func GoIdent() string {
	name := String(-1)
	for {
		firstLetter := rune(name[0])
		if !unicode.IsDigit(firstLetter) {
			break
		}
		name = String(RandLen)
	}
	return name
}
