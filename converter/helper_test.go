package converter

import (
	"math"
	"reflect"
	"slices"

	"github.com/naivary/codemark/maker"
	"github.com/naivary/codemark/registry"
	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var defsSet = newDefsSet()

var mngr = newManager()

func newManager() sdk.ConverterManager {
	mngr, err := NewManager(defsSet)
	if err != nil {
		panic(err)
	}
	return mngr
}

func newDefsSet() sdk.Registry {
	reg, err := sdktesting.NewDefsSet(registry.InMemory(), maker.New())
	if err != nil {
		panic(err)
	}
	return reg
}

func newConvTester(conv sdk.Converter) (sdktesting.ConverterTester, error) {
	tester, err := sdktesting.NewConvTester(conv, fromMap(conv))
	if err != nil {
		return nil, err
	}
	for _, typ := range toTypes(conv) {
		to := reflect.TypeOf(typ)
		vvfn := getVVFn(to)
		if err := tester.AddVVFunc(to, vvfn); err != nil {
			return nil, err
		}
	}
	return tester, nil
}

func fromMap(conv sdk.Converter) map[reflect.Type]reflect.Type {
	if _, isList := conv.(*listConverter); isList {
		return fromMapOfList()
	}
	if _, isInt := conv.(*intConverter); isInt {
		return fromMapOfInteger()
	}
	if _, isFloat := conv.(*floatConverter); isFloat {
		return fromMapOfFloat()
	}
	if _, isComplex := conv.(*complexConverter); isComplex {
		return fromMapOfComplex()
	}
	if _, isBool := conv.(*boolConverter); isBool {
		return fromMapOfBool()
	}
	if _, isString := conv.(*stringConverter); isString {
		return fromMapOfString()
	}
	return nil

}

func fromMapOfList() map[reflect.Type]reflect.Type {
	return map[reflect.Type]reflect.Type{
		// Signed integer slices
		reflect.TypeFor[[]int]():   reflect.TypeFor[sdktesting.IntList](),
		reflect.TypeFor[[]int8]():  reflect.TypeFor[sdktesting.I8List](),
		reflect.TypeFor[[]int16](): reflect.TypeFor[sdktesting.I16List](),
		reflect.TypeFor[[]int32](): reflect.TypeFor[sdktesting.I32List](),
		reflect.TypeFor[[]int64](): reflect.TypeFor[sdktesting.I64List](),
		// Unsigned integer slices
		reflect.TypeFor[[]uint]():   reflect.TypeFor[sdktesting.UintList](),
		reflect.TypeFor[[]uint8]():  reflect.TypeFor[sdktesting.U8List](),
		reflect.TypeFor[[]uint16](): reflect.TypeFor[sdktesting.U16List](),
		reflect.TypeFor[[]uint32](): reflect.TypeFor[sdktesting.U32List](),
		reflect.TypeFor[[]uint64](): reflect.TypeFor[sdktesting.U64List](),
		// Float slices
		reflect.TypeFor[[]float32](): reflect.TypeFor[sdktesting.F32List](),
		reflect.TypeFor[[]float64](): reflect.TypeFor[sdktesting.F64List](),
		// Complex slices
		reflect.TypeFor[[]complex64]():  reflect.TypeFor[sdktesting.C64List](),
		reflect.TypeFor[[]complex128](): reflect.TypeFor[sdktesting.C128List](),
		// String and bool slices
		reflect.TypeFor[[]string](): reflect.TypeFor[sdktesting.StringList](),
		reflect.TypeFor[[]bool]():   reflect.TypeFor[sdktesting.BoolList](),
		// Pointer to signed integer slices
		reflect.TypeFor[[]*int]():   reflect.TypeFor[sdktesting.PtrIntList](),
		reflect.TypeFor[[]*int8]():  reflect.TypeFor[sdktesting.PtrI8List](),
		reflect.TypeFor[[]*int16](): reflect.TypeFor[sdktesting.PtrI16List](),
		reflect.TypeFor[[]*int32](): reflect.TypeFor[sdktesting.PtrI32List](),
		reflect.TypeFor[[]*int64](): reflect.TypeFor[sdktesting.PtrI64List](),
		// Pointer to unsigned integer slices
		reflect.TypeFor[[]*uint]():   reflect.TypeFor[sdktesting.PtrUintList](),
		reflect.TypeFor[[]*uint8]():  reflect.TypeFor[sdktesting.PtrU8List](),
		reflect.TypeFor[[]*uint16](): reflect.TypeFor[sdktesting.PtrU16List](),
		reflect.TypeFor[[]*uint32](): reflect.TypeFor[sdktesting.PtrU32List](),
		reflect.TypeFor[[]*uint64](): reflect.TypeFor[sdktesting.PtrU64List](),
		// Pointer to float slices
		reflect.TypeFor[[]*float32](): reflect.TypeFor[sdktesting.PtrF32List](),
		reflect.TypeFor[[]*float64](): reflect.TypeFor[sdktesting.PtrF64List](),
		// Pointer to complex slices
		reflect.TypeFor[[]*complex64]():  reflect.TypeFor[sdktesting.PtrC64List](),
		reflect.TypeFor[[]*complex128](): reflect.TypeFor[sdktesting.PtrC128List](),
		// Pointer to string and bool slices
		reflect.TypeFor[[]*string](): reflect.TypeFor[sdktesting.PtrStringList](),
		reflect.TypeFor[[]*bool]():   reflect.TypeFor[sdktesting.PtrBoolList](),
	}
}

func fromMapOfInteger() map[reflect.Type]reflect.Type {
	return map[reflect.Type]reflect.Type{
		// Signed integers
		reflect.TypeFor[int]():    reflect.TypeFor[sdktesting.Int](),
		reflect.TypeFor[int8]():   reflect.TypeFor[sdktesting.I8](),
		reflect.TypeFor[int16]():  reflect.TypeFor[sdktesting.I16](),
		reflect.TypeFor[int32]():  reflect.TypeFor[sdktesting.I32](),
		reflect.TypeFor[int64]():  reflect.TypeFor[sdktesting.I64](),
		reflect.TypeFor[*int]():   reflect.TypeFor[sdktesting.PtrInt](),
		reflect.TypeFor[*int8]():  reflect.TypeFor[sdktesting.PtrI8](),
		reflect.TypeFor[*int16](): reflect.TypeFor[sdktesting.PtrI16](),
		reflect.TypeFor[*int32](): reflect.TypeFor[sdktesting.PtrI32](),
		reflect.TypeFor[*int64](): reflect.TypeFor[sdktesting.PtrI64](),
		// Unsigned integers
		reflect.TypeFor[uint]():    reflect.TypeFor[sdktesting.Uint](),
		reflect.TypeFor[uint8]():   reflect.TypeFor[sdktesting.U8](),
		reflect.TypeFor[uint16]():  reflect.TypeFor[sdktesting.U16](),
		reflect.TypeFor[uint32]():  reflect.TypeFor[sdktesting.U32](),
		reflect.TypeFor[uint64]():  reflect.TypeFor[sdktesting.U64](),
		reflect.TypeFor[*uint]():   reflect.TypeFor[sdktesting.PtrUint](),
		reflect.TypeFor[*uint8]():  reflect.TypeFor[sdktesting.PtrU8](),
		reflect.TypeFor[*uint16](): reflect.TypeFor[sdktesting.PtrU16](),
		reflect.TypeFor[*uint32](): reflect.TypeFor[sdktesting.PtrU32](),
		reflect.TypeFor[*uint64](): reflect.TypeFor[sdktesting.PtrU64](),
	}
}

func fromMapOfComplex() map[reflect.Type]reflect.Type {
	return map[reflect.Type]reflect.Type{
		reflect.TypeFor[complex64]():   reflect.TypeFor[sdktesting.C64](),
		reflect.TypeFor[complex128]():  reflect.TypeFor[sdktesting.C128](),
		reflect.TypeFor[*complex64]():  reflect.TypeFor[sdktesting.PtrC64](),
		reflect.TypeFor[*complex128](): reflect.TypeFor[sdktesting.PtrC128](),
	}
}

func fromMapOfString() map[reflect.Type]reflect.Type {
	return map[reflect.Type]reflect.Type{
		reflect.TypeFor[string]():  reflect.TypeFor[sdktesting.String](),
		reflect.TypeFor[*string](): reflect.TypeFor[sdktesting.PtrString](),
	}
}

func fromMapOfFloat() map[reflect.Type]reflect.Type {
	return map[reflect.Type]reflect.Type{
		reflect.TypeFor[float32]():  reflect.TypeFor[sdktesting.F32](),
		reflect.TypeFor[float64]():  reflect.TypeFor[sdktesting.F64](),
		reflect.TypeFor[*float32](): reflect.TypeFor[sdktesting.PtrF32](),
		reflect.TypeFor[*float64](): reflect.TypeFor[sdktesting.PtrF64](),
	}
}

func fromMapOfBool() map[reflect.Type]reflect.Type {
	return map[reflect.Type]reflect.Type{
		reflect.TypeFor[bool]():  reflect.TypeFor[sdktesting.Bool](),
		reflect.TypeFor[*bool](): reflect.TypeFor[sdktesting.PtrBool](),
	}
}

func toTypes(conv sdk.Converter) []any {
	if _, isList := conv.(*listConverter); isList {
		return sdktesting.ListTypes()
	}
	if _, isInt := conv.(*intConverter); isInt {
		return slices.Concat(sdktesting.IntTypes(), sdktesting.UintTypes())
	}
	if _, isFloat := conv.(*floatConverter); isFloat {
		return sdktesting.FloatTypes()
	}
	if _, isComplex := conv.(*complexConverter); isComplex {
		return sdktesting.ComplexTypes()
	}
	if _, isBool := conv.(*boolConverter); isBool {
		return sdktesting.BoolTypes()
	}
	if _, isString := conv.(*stringConverter); isString {
		return sdktesting.StringTypes()
	}
	return nil

}

func getVVFn(rtype reflect.Type) sdktesting.ValidValueFunc {
	if sdkutil.IsValidSlice(rtype) {
		return isValidList
	}
	if sdkutil.IsInt(rtype) || sdkutil.IsUint(rtype) {
		return isValidInteger
	}
	if sdkutil.IsFloat(rtype) {
		return isValidFloat
	}
	if sdkutil.IsComplex(rtype) {
		return isValidComplex
	}
	if sdkutil.IsBool(rtype) {
		return isValidBool
	}
	if sdkutil.IsString(rtype) {
		return isValidString
	}
	return nil
}

func isValidInteger(got, want reflect.Value) bool {
	got = sdkutil.DeRefValue(got)
	var i64 int64
	switch got.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 = got.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64 := got.Uint()
		if u64 > math.MaxInt64 {
			return false // uint64 value too large for int64
		}
		i64 = int64(u64)
	default:
		return false
	}
	w := want.Interface().(int64)
	return i64 == w
}

func isValidList(got, want reflect.Value) bool {
	elem := got.Type().Elem()
	vvfn := getVVFn(elem)
	if vvfn == nil {
		return false
	}
	for i := range want.Len() {
		wantElem := want.Index(i)
		gotElem := got.Index(i)
		if !vvfn(gotElem, wantElem) {
			return false
		}
	}
	return true
}

func isValidFloat(got, want reflect.Value) bool {
	got = sdkutil.DeRefValue(got)
	w := want.Interface().(float64)
	return sdktesting.AlmostEqual(got.Float(), w)
}

func isValidBool(got, want reflect.Value) bool {
	got = sdkutil.DeRefValue(got)
	w := want.Interface().(bool)
	return got.Bool() == w
}

func isValidString(got, want reflect.Value) bool {
	got = sdkutil.DeRefValue(got)
	w := want.Interface().(string)
	return got.String() == w
}

func isValidComplex(got, want reflect.Value) bool {
	got = sdkutil.DeRefValue(got)
	w := want.Interface().(complex128)
	return got.Complex() == w
}
