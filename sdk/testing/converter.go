package testing

import (
	"fmt"
	"math"
	"reflect"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
	"github.com/naivary/codemark/sdk/utils"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

type ConverterTestCase struct {
	Name         string
	Marker       parser.Marker
	Target       sdk.Target
	ToType       reflect.Type
	IsValidCase  bool
	IsValidValue func(got reflect.Value, want reflect.Value) bool
}

type ConverterTester interface {
	NewTest(conv sdk.Converter) ([]ConverterTestCase, error)
}

type converterTester struct {
	vvfns map[string]func(got, want reflect.Value) bool
	types map[string]reflect.Type
}

func NewConverterTester(vvfns map[string]func(got, want reflect.Value) bool, types map[string]reflect.Type) (ConverterTester, error) {
	c := &converterTester{}
	c.defaultVVFns()
	c.defaultTypes()
	for typeID, fn := range vvfns {
		_, found := c.vvfns[typeID]
		if found {
			return nil, fmt.Errorf("IsValidFunction exists: %s\n", typeID)
		}
		c.vvfns[typeID] = fn
	}
	for typeID, rtype := range types {
		_, found := c.types[typeID]
		if found {
			return nil, fmt.Errorf("type id already exists: %s\n", typeID)
		}
		c.types[typeID] = rtype
	}
	return c, nil
}

func (c *converterTester) NewTest(conv sdk.Converter) ([]ConverterTestCase, error) {
	tests := make([]ConverterTestCase, 0, len(conv.SupportedTypes()))
	for _, rtype := range conv.SupportedTypes() {
		typeID := sdkutil.TypeID(rtype)
		marker := RandMarkerFromRefType(rtype)
		if marker == nil {
			return nil, fmt.Errorf("no valid marker found: %v\n", rtype)
		}
		tc := ConverterTestCase{
			Name:         "random-name",
			Marker:       marker,
			Target:       sdk.TargetAny,
			ToType:       c.types[typeID],
			IsValidCase:  true,
			IsValidValue: c.vvfns[typeID],
		}
		tests = append(tests, tc)
	}
	return tests, nil
}

func (c *converterTester) defaultTypes() {
	types := DefaultTypes()
	c.types = make(map[string]reflect.Type, len(types))
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		typeID := sdkutil.TypeID(rtype)
		c.types[typeID] = rtype
	}
}

func (c *converterTester) defaultVVFns() {
	c.vvfns = map[string]func(got, want reflect.Value) bool{
		// Integers
		sdkutil.TypeIDFromAny(Int(0)):       isValidInteger,
		sdkutil.TypeIDFromAny(I8(0)):        isValidInteger,
		sdkutil.TypeIDFromAny(I16(0)):       isValidInteger,
		sdkutil.TypeIDFromAny(Byte(0)):      isValidInteger,
		sdkutil.TypeIDFromAny(I32(0)):       isValidInteger,
		sdkutil.TypeIDFromAny(I64(0)):       isValidInteger,
		sdkutil.TypeIDFromAny(PtrInt(nil)):  isValidInteger,
		sdkutil.TypeIDFromAny(PtrI8(nil)):   isValidInteger,
		sdkutil.TypeIDFromAny(PtrI16(nil)):  isValidInteger,
		sdkutil.TypeIDFromAny(PtrByte(nil)): isValidInteger,
		sdkutil.TypeIDFromAny(PtrI32(nil)):  isValidInteger,
		sdkutil.TypeIDFromAny(PtrI64(nil)):  isValidInteger,

		// Unsigned integers
		sdkutil.TypeIDFromAny(Uint(0)):      isValidInteger,
		sdkutil.TypeIDFromAny(Rune(0)):      isValidInteger,
		sdkutil.TypeIDFromAny(U8(0)):        isValidInteger,
		sdkutil.TypeIDFromAny(U16(0)):       isValidInteger,
		sdkutil.TypeIDFromAny(U32(0)):       isValidInteger,
		sdkutil.TypeIDFromAny(U64(0)):       isValidInteger,
		sdkutil.TypeIDFromAny(PtrUint(nil)): isValidInteger,
		sdkutil.TypeIDFromAny(PtrRune(nil)): isValidInteger,
		sdkutil.TypeIDFromAny(PtrU8(nil)):   isValidInteger,
		sdkutil.TypeIDFromAny(PtrU16(nil)):  isValidInteger,
		sdkutil.TypeIDFromAny(PtrU32(nil)):  isValidInteger,
		sdkutil.TypeIDFromAny(PtrU64(nil)):  isValidInteger,

		// Floats
		sdkutil.TypeIDFromAny(F32(0)):      isValidFloat,
		sdkutil.TypeIDFromAny(F64(0)):      isValidFloat,
		sdkutil.TypeIDFromAny(PtrF32(nil)): isValidFloat,
		sdkutil.TypeIDFromAny(PtrF64(nil)): isValidFloat,

		// Complex
		sdkutil.TypeIDFromAny(C64(0)):       isValidComplex,
		sdkutil.TypeIDFromAny(C128(0)):      isValidComplex,
		sdkutil.TypeIDFromAny(PtrC64(nil)):  isValidComplex,
		sdkutil.TypeIDFromAny(PtrC128(nil)): isValidComplex,

		// Booleans and strings
		sdkutil.TypeIDFromAny(Bool(false)):    isValidBool,
		sdkutil.TypeIDFromAny(String("")):     isValidString,
		sdkutil.TypeIDFromAny(PtrBool(nil)):   isValidBool,
		sdkutil.TypeIDFromAny(PtrString(nil)): isValidString,

		// List
		sdkutil.TypeIDFromAny(IntList(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(I8List(nil)):     c.isValidList,
		sdkutil.TypeIDFromAny(I16List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(ByteList(nil)):   c.isValidList,
		sdkutil.TypeIDFromAny(I64List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(UintList(nil)):   c.isValidList,
		sdkutil.TypeIDFromAny(RuneList(nil)):   c.isValidList,
		sdkutil.TypeIDFromAny(U16List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(U32List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(U64List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(F32List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(F64List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(C64List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(C128List(nil)):   c.isValidList,
		sdkutil.TypeIDFromAny(BoolList(nil)):   c.isValidList,
		sdkutil.TypeIDFromAny(StringList(nil)): c.isValidList,

		// List of pointers
		sdkutil.TypeIDFromAny(PtrIntList(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(PtrI8List(nil)):     c.isValidList,
		sdkutil.TypeIDFromAny(PtrI16List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(PtrByteList(nil)):   c.isValidList,
		sdkutil.TypeIDFromAny(PtrI64List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(PtrUintList(nil)):   c.isValidList,
		sdkutil.TypeIDFromAny(PtrRuneList(nil)):   c.isValidList,
		sdkutil.TypeIDFromAny(PtrU16List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(PtrU32List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(PtrU64List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(PtrF32List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(PtrF64List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(PtrC64List(nil)):    c.isValidList,
		sdkutil.TypeIDFromAny(PtrC128List(nil)):   c.isValidList,
		sdkutil.TypeIDFromAny(PtrBoolList(nil)):   c.isValidList,
		sdkutil.TypeIDFromAny(PtrStringList(nil)): c.isValidList,
	}
}

func (c *converterTester) isValidList(got, want reflect.Value) bool {
	elem := got.Elem()
	vvfn, found := c.vvfns[sdkutil.TypeID(elem.Type())]
	if !found {
		return false
	}
	i := 0
	for wantElem := range want.Seq() {
		gotElem := got.Index(i)
		if !vvfn(gotElem, wantElem) {
			return false
		}
	}
	return true
}

func isValidInteger(got, want reflect.Value) bool {
	if utils.IsPointer(got.Type()) {
		got = got.Elem()
	}
	var i64 int64
	switch got.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 = got.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64 := got.Uint()
		if u64 > math.MaxInt64 {
			return false // uint64 value too large for int64
		}
	default:
		return false
	}
	w := want.Interface().(int64)
	return i64 == w
}

func isValidFloat(got, want reflect.Value) bool {
	if utils.IsPointer(got.Type()) {
		got = got.Elem()
	}
	w := want.Interface().(float64)
	return AlmostEqual(got.Float(), w)
}

func isValidComplex(got, want reflect.Value) bool {
	if utils.IsPointer(got.Type()) {
		got = got.Elem()
	}
	w := want.Interface().(complex128)
	return got.Complex() == w
}

func isValidString(got, want reflect.Value) bool {
	if utils.IsPointer(got.Type()) {
		got = got.Elem()
	}
	w := want.Interface().(string)
	return got.String() == w
}

func isValidBool(got, want reflect.Value) bool {
	if utils.IsPointer(got.Type()) {
		got = got.Elem()
	}
	w := want.Interface().(bool)
	return got.Bool() == w
}
