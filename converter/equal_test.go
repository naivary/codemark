package converter

import (
	"reflect"
	"time"

	"github.com/naivary/codemark/internal/equal"
	"github.com/naivary/codemark/registry/registrytest"
	"github.com/naivary/codemark/rtypeutil"
)

var equalFuncs = map[reflect.Type]func(got, want reflect.Value) bool{
	reflect.TypeFor[Duration]():             equalDuration,
	reflect.TypeFor[PtrDuration]():          equalDuration,
	reflect.TypeFor[Time]():                 equalTime,
	reflect.TypeFor[PtrTime]():              equalTime,
	reflect.TypeFor[registrytest.Byte]():    equalByte,
	reflect.TypeFor[registrytest.PtrByte](): equalByte,
	reflect.TypeFor[registrytest.Rune]():    equalRune,
	reflect.TypeFor[registrytest.PtrRune](): equalRune,
}

func getEqualFunc(t reflect.Type) func(got, want reflect.Value) bool {
	fn, ok := equalFuncs[t]
	if !ok {
		return equal.GetFunc(t)
	}
	return fn
}

func equalTime(got, want reflect.Value) bool {
	t := reflect.TypeOf(time.Time{})
	gotTime := rtypeutil.DerefValue(got).Convert(t).Interface().(time.Time)
	wantTime, err := time.Parse(time.RFC3339, want.String())
	if err != nil {
		return false
	}
	return gotTime.Equal(wantTime)
}

func equalDuration(got, want reflect.Value) bool {
	t := reflect.TypeFor[time.Duration]()
	gotDuration := rtypeutil.DerefValue(got).Convert(t).Interface().(time.Duration)
	wantDuration, err := time.ParseDuration(want.String())
	if err != nil {
		return false
	}
	return gotDuration == wantDuration
}

func equalByte(got, want reflect.Value) bool {
	got = rtypeutil.DerefValue(got)
	return want.String() == string(byte(got.Uint()))
}

func equalRune(got, want reflect.Value) bool {
	got = rtypeutil.DerefValue(got)
	return want.String() == string(rune(got.Int()))
}
