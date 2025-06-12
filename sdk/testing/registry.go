package testing

import (
	"reflect"
	"slices"

	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

var typeIDToMarkerName = map[string]string{
	sdkutil.TypeID(reflect.TypeOf(String(""))):         "path:to:str",
	sdkutil.TypeID(reflect.TypeOf(PtrString(nil))):     "path:to:ptrstr",
	sdkutil.TypeID(reflect.TypeOf(Bool(false))):        "path:to:bool",
	sdkutil.TypeID(reflect.TypeOf(PtrBool(nil))):       "path:to:ptrbool",
	sdkutil.TypeID(reflect.TypeOf(C64(0 + 0i))):        "path:to:c64",
	sdkutil.TypeID(reflect.TypeOf(C128(0 + 0i))):       "path:to:c128",
	sdkutil.TypeID(reflect.TypeOf(PtrC64(nil))):        "path:to:ptrc64",
	sdkutil.TypeID(reflect.TypeOf(PtrC128(nil))):       "path:to:ptrc128",
	sdkutil.TypeID(reflect.TypeOf(Int(0))):             "path:to:i",
	sdkutil.TypeID(reflect.TypeOf(I8(0))):              "path:to:i8",
	sdkutil.TypeID(reflect.TypeOf(I16(0))):             "path:to:i16",
	sdkutil.TypeID(reflect.TypeOf(I32(0))):             "path:to:i32",
	sdkutil.TypeID(reflect.TypeOf(I64(0))):             "path:to:i64",
	sdkutil.TypeID(reflect.TypeOf(PtrInt(nil))):        "path:to:ptri",
	sdkutil.TypeID(reflect.TypeOf(PtrI8(nil))):         "path:to:ptri8",
	sdkutil.TypeID(reflect.TypeOf(PtrI16(nil))):        "path:to:ptri16",
	sdkutil.TypeID(reflect.TypeOf(PtrI32(nil))):        "path:to:ptri32",
	sdkutil.TypeID(reflect.TypeOf(PtrI64(nil))):        "path:to:ptri64",
	sdkutil.TypeID(reflect.TypeOf(Uint(0))):            "path:to:ui",
	sdkutil.TypeID(reflect.TypeOf(U8(0))):              "path:to:ui8",
	sdkutil.TypeID(reflect.TypeOf(U16(0))):             "path:to:ui16",
	sdkutil.TypeID(reflect.TypeOf(U32(0))):             "path:to:ui32",
	sdkutil.TypeID(reflect.TypeOf(U64(0))):             "path:to:ui64",
	sdkutil.TypeID(reflect.TypeOf(PtrUint(nil))):       "path:to:ptrui",
	sdkutil.TypeID(reflect.TypeOf(PtrU8(nil))):         "path:to:ptrui8",
	sdkutil.TypeID(reflect.TypeOf(PtrU16(nil))):        "path:to:ptrui16",
	sdkutil.TypeID(reflect.TypeOf(PtrU32(nil))):        "path:to:ptrui32",
	sdkutil.TypeID(reflect.TypeOf(PtrU64(nil))):        "path:to:ptrui64",
	sdkutil.TypeID(reflect.TypeOf(Byte(0))):            "path:to:byte",
	sdkutil.TypeID(reflect.TypeOf(Rune(0))):            "path:to:rune",
	sdkutil.TypeID(reflect.TypeOf(PtrByte(nil))):       "path:to:ptrbyte",
	sdkutil.TypeID(reflect.TypeOf(PtrRune(nil))):       "path:to:ptrrune",
	sdkutil.TypeID(reflect.TypeOf(F32(0.0))):           "path:to:f32",
	sdkutil.TypeID(reflect.TypeOf(F64(0.0))):           "path:to:f64",
	sdkutil.TypeID(reflect.TypeOf(PtrF32(nil))):        "path:to:ptrf32",
	sdkutil.TypeID(reflect.TypeOf(PtrF64(nil))):        "path:to:ptrf64",
	sdkutil.TypeID(reflect.TypeOf(StringList(nil))):    "path:to:stringlist",
	sdkutil.TypeID(reflect.TypeOf(IntList(nil))):       "path:to:intlist",
	sdkutil.TypeID(reflect.TypeOf(I8List(nil))):        "path:to:i8list",
	sdkutil.TypeID(reflect.TypeOf(I16List(nil))):       "path:to:i16list",
	sdkutil.TypeID(reflect.TypeOf(ByteList(nil))):      "path:to:bytelist",
	sdkutil.TypeID(reflect.TypeOf(I64List(nil))):       "path:to:i64list",
	sdkutil.TypeID(reflect.TypeOf(UintList(nil))):      "path:to:uintlist",
	sdkutil.TypeID(reflect.TypeOf(RuneList(nil))):      "path:to:runelist",
	sdkutil.TypeID(reflect.TypeOf(U16List(nil))):       "path:to:ui16list",
	sdkutil.TypeID(reflect.TypeOf(U32List(nil))):       "path:to:ui32list",
	sdkutil.TypeID(reflect.TypeOf(U64List(nil))):       "path:to:ui64list",
	sdkutil.TypeID(reflect.TypeOf(F32List(nil))):       "path:to:f32list",
	sdkutil.TypeID(reflect.TypeOf(F64List(nil))):       "path:to:f64list",
	sdkutil.TypeID(reflect.TypeOf(C64List(nil))):       "path:to:c64list",
	sdkutil.TypeID(reflect.TypeOf(C128List(nil))):      "path:to:c128list",
	sdkutil.TypeID(reflect.TypeOf(BoolList(nil))):      "path:to:boollist",
	sdkutil.TypeID(reflect.TypeOf(PtrStringList(nil))): "path:to:ptrstringlist",
	sdkutil.TypeID(reflect.TypeOf(PtrIntList(nil))):    "path:to:ptrintlist",
	sdkutil.TypeID(reflect.TypeOf(PtrI8List(nil))):     "path:to:ptri8list",
	sdkutil.TypeID(reflect.TypeOf(PtrI16List(nil))):    "path:to:ptri16list",
	sdkutil.TypeID(reflect.TypeOf(PtrByteList(nil))):   "path:to:ptrbytelist",
	sdkutil.TypeID(reflect.TypeOf(PtrI64List(nil))):    "path:to:ptri64list",
	sdkutil.TypeID(reflect.TypeOf(PtrUintList(nil))):   "path:to:ptruintlist",
	sdkutil.TypeID(reflect.TypeOf(PtrRuneList(nil))):   "path:to:ptrrunelist",
	sdkutil.TypeID(reflect.TypeOf(PtrU16List(nil))):    "path:to:ptrui16list",
	sdkutil.TypeID(reflect.TypeOf(PtrU32List(nil))):    "path:to:ptrui32list",
	sdkutil.TypeID(reflect.TypeOf(PtrU64List(nil))):    "path:to:ptrui64list",
	sdkutil.TypeID(reflect.TypeOf(PtrF32List(nil))):    "path:to:ptrf32list",
	sdkutil.TypeID(reflect.TypeOf(PtrF64List(nil))):    "path:to:ptrf64list",
	sdkutil.TypeID(reflect.TypeOf(PtrC64List(nil))):    "path:to:ptrc64list",
	sdkutil.TypeID(reflect.TypeOf(PtrC128List(nil))):   "path:to:ptrc128list",
	sdkutil.TypeID(reflect.TypeOf(PtrBoolList(nil))):   "path:to:ptrboollist",
}

type F32 float32
type F64 float64

type PtrF32 *float32
type PtrF64 *float64

type Bool bool
type PtrBool *bool

type C64 complex64
type C128 complex128

type PtrC64 *complex64
type PtrC128 *complex128

type String string
type PtrString *string

type Int int
type I8 int8
type I16 int16

// byte=int32
type Byte byte
type I32 int32
type I64 int64

type PtrInt *int
type PtrI8 *int8
type PtrI16 *int16

// *byte=*int32
type PtrByte *byte
type PtrI32 *int32
type PtrI64 *int64

type Uint uint

// rune=uint8
type Rune rune
type U8 uint8
type U16 uint16
type U32 uint32
type U64 uint64

type PtrUint *uint
type PtrU8 *uint8

// *rune=*uint8
type PtrRune *rune
type PtrU16 *uint16
type PtrU32 *uint32
type PtrU64 *uint64

// string list
type StringList []string

// int list
type IntList []int
type I8List []int8
type I16List []int16
type ByteList []byte // int32
type I64List []int64

// uint list
type UintList []uint
type RuneList []rune // uint8
type U16List []uint16
type U32List []uint32
type U64List []uint64

// float list
type F32List []float32
type F64List []float64

// complex list
type C64List []complex64
type C128List []complex128

// bool list
type BoolList []bool

// ptr string list
type PtrStringList []*string

// ptr bool list
type PtrBoolList []*bool

// ptr int list
type PtrIntList []*int
type PtrI8List []*int8
type PtrI16List []*int16
type PtrByteList []*byte // int32
type PtrI64List []*int64

// ptr uint list
type PtrUintList []*uint
type PtrRuneList []*rune // uint8
type PtrU16List []*uint16
type PtrU32List []*uint32
type PtrU64List []*uint64

// ptr float
type PtrF32List []*float32
type PtrF64List []*float64

// complex
type PtrC64List []*complex64
type PtrC128List []*complex128

func NewDefsSet(reg sdk.Registry, b sdk.DefinitionMaker, customDefs ...*sdk.Definition) (sdk.Registry, error) {
	defs := []*sdk.Definition{
		// string
		b.MustMakeDef("path:to:str", reflect.TypeOf(String("")), sdk.TargetAny),
		// ptr string
		b.MustMakeDef("path:to:ptrstr", reflect.TypeOf(PtrString(new(string))), sdk.TargetAny),
		// bool
		b.MustMakeDef("path:to:bool", reflect.TypeOf(Bool(false)), sdk.TargetAny),
		// ptr bool
		b.MustMakeDef("path:to:ptrbool", reflect.TypeOf(PtrBool(new(bool))), sdk.TargetAny),
		// complex
		b.MustMakeDef("path:to:c64", reflect.TypeOf(C64(0+0i)), sdk.TargetAny),
		b.MustMakeDef("path:to:c128", reflect.TypeOf(C128(0+0i)), sdk.TargetAny),
		// ptr complex
		b.MustMakeDef("path:to:ptrc64", reflect.TypeOf(PtrC64(new(complex64))), sdk.TargetAny),
		b.MustMakeDef("path:to:ptrc128", reflect.TypeOf(PtrC128(new(complex128))), sdk.TargetAny),
		// int
		b.MustMakeDef("path:to:i", reflect.TypeOf(Int(0)), sdk.TargetAny),
		b.MustMakeDef("path:to:i8", reflect.TypeOf(I8(0)), sdk.TargetAny),
		b.MustMakeDef("path:to:i16", reflect.TypeOf(I16(0)), sdk.TargetAny),
		b.MustMakeDef("path:to:i32", reflect.TypeOf(I32(0)), sdk.TargetAny),
		b.MustMakeDef("path:to:i64", reflect.TypeOf(I64(0)), sdk.TargetAny),
		// ptr int
		b.MustMakeDef("path:to:ptri", reflect.TypeOf(PtrInt(new(int))), sdk.TargetAny),
		b.MustMakeDef("path:to:ptri8", reflect.TypeOf(PtrI8(new(int8))), sdk.TargetAny),
		b.MustMakeDef("path:to:ptri16", reflect.TypeOf(PtrI16(new(int16))), sdk.TargetAny),
		b.MustMakeDef("path:to:ptri32", reflect.TypeOf(PtrI32(new(int32))), sdk.TargetAny),
		b.MustMakeDef("path:to:ptri64", reflect.TypeOf(PtrI64(new(int64))), sdk.TargetAny),
		// uint
		b.MustMakeDef("path:to:ui", reflect.TypeOf(Uint(0)), sdk.TargetAny),
		b.MustMakeDef("path:to:ui8", reflect.TypeOf(U8(0)), sdk.TargetAny),
		b.MustMakeDef("path:to:ui16", reflect.TypeOf(U16(0)), sdk.TargetAny),
		b.MustMakeDef("path:to:ui32", reflect.TypeOf(U32(0)), sdk.TargetAny),
		b.MustMakeDef("path:to:ui64", reflect.TypeOf(U64(0)), sdk.TargetAny),
		// ptr uint
		b.MustMakeDef("path:to:ptrui", reflect.TypeOf(PtrUint(new(uint))), sdk.TargetAny),
		b.MustMakeDef("path:to:ptrui8", reflect.TypeOf(PtrU8(new(uint8))), sdk.TargetAny),
		b.MustMakeDef("path:to:ptrui16", reflect.TypeOf(PtrU16(new(uint16))), sdk.TargetAny),
		b.MustMakeDef("path:to:ptrui32", reflect.TypeOf(PtrU32(new(uint32))), sdk.TargetAny),
		b.MustMakeDef("path:to:ptrui64", reflect.TypeOf(PtrU64(new(uint64))), sdk.TargetAny),
		// byte and rune
		b.MustMakeDef("path:to:byte", reflect.TypeOf(Byte(0)), sdk.TargetAny),
		b.MustMakeDef("path:to:rune", reflect.TypeOf(Rune(0)), sdk.TargetAny),
		// ptr byte and rune
		b.MustMakeDef("path:to:ptrbyte", reflect.TypeOf(PtrByte(new(byte))), sdk.TargetAny),
		b.MustMakeDef("path:to:ptrrune", reflect.TypeOf(PtrRune(new(rune))), sdk.TargetAny),
		// float
		b.MustMakeDef("path:to:f32", reflect.TypeOf(F32(0.0)), sdk.TargetAny),
		b.MustMakeDef("path:to:f64", reflect.TypeOf(F64(0.0)), sdk.TargetAny),
		// ptr float
		b.MustMakeDef("path:to:ptrf32", reflect.TypeOf(PtrF32(new(float32))), sdk.TargetAny),
		b.MustMakeDef("path:to:ptrf64", reflect.TypeOf(PtrF64(new(float64))), sdk.TargetAny),
		// string list
		b.MustMakeDef("path:to:stringlist", reflect.TypeOf(StringList([]string{})), sdk.TargetAny),
		// int list
		b.MustMakeDef("path:to:intlist", reflect.TypeOf(IntList([]int{})), sdk.TargetAny),
		b.MustMakeDef("path:to:i8list", reflect.TypeOf(I8List([]int8{})), sdk.TargetAny),
		b.MustMakeDef("path:to:i16list", reflect.TypeOf(I16List([]int16{})), sdk.TargetAny),
		b.MustMakeDef("path:to:bytelist", reflect.TypeOf(ByteList([]byte{})), sdk.TargetAny),
		b.MustMakeDef("path:to:i64list", reflect.TypeOf(I64List([]int64{})), sdk.TargetAny),
		// uint list
		b.MustMakeDef("path:to:uintlist", reflect.TypeOf(UintList([]uint{})), sdk.TargetAny),
		b.MustMakeDef("path:to:runelist", reflect.TypeOf(RuneList([]rune{})), sdk.TargetAny),
		b.MustMakeDef("path:to:ui16list", reflect.TypeOf(U16List([]uint16{})), sdk.TargetAny),
		b.MustMakeDef("path:to:ui32list", reflect.TypeOf(U32List([]uint32{})), sdk.TargetAny),
		b.MustMakeDef("path:to:ui64list", reflect.TypeOf(U64List([]uint64{})), sdk.TargetAny),
		// float list
		b.MustMakeDef("path:to:f32list", reflect.TypeOf(F32List([]float32{})), sdk.TargetAny),
		b.MustMakeDef("path:to:f64list", reflect.TypeOf(F64List([]float64{})), sdk.TargetAny),
		// complex list
		b.MustMakeDef("path:to:c64list", reflect.TypeOf(C64List([]complex64{})), sdk.TargetAny),
		b.MustMakeDef("path:to:c128list", reflect.TypeOf(C128List([]complex128{})), sdk.TargetAny),
		// bool list
		b.MustMakeDef("path:to:boollist", reflect.TypeOf(BoolList([]bool{})), sdk.TargetAny),
		// ptr string list
		b.MustMakeDef("path:to:ptrstringlist", reflect.TypeOf(PtrStringList([]*string{})), sdk.TargetAny),
		// ptr int list
		b.MustMakeDef("path:to:ptrintlist", reflect.TypeOf(PtrIntList([]*int{})), sdk.TargetAny),
		b.MustMakeDef("path:to:ptri8list", reflect.TypeOf(PtrI8List([]*int8{})), sdk.TargetAny),
		b.MustMakeDef("path:to:ptri16list", reflect.TypeOf(PtrI16List([]*int16{})), sdk.TargetAny),
		b.MustMakeDef("path:to:ptrbytelist", reflect.TypeOf(PtrByteList([]*byte{})), sdk.TargetAny),
		b.MustMakeDef("path:to:ptri64list", reflect.TypeOf(PtrI64List([]*int64{})), sdk.TargetAny),
		// ptr uint list
		b.MustMakeDef("path:to:ptruintlist", reflect.TypeOf(PtrUintList([]*uint{})), sdk.TargetAny),
		b.MustMakeDef("path:to:ptrrunelist", reflect.TypeOf(PtrRuneList([]*rune{})), sdk.TargetAny),
		b.MustMakeDef("path:to:ptrui16list", reflect.TypeOf(PtrU16List([]*uint16{})), sdk.TargetAny),
		b.MustMakeDef("path:to:ptrui32list", reflect.TypeOf(PtrU32List([]*uint32{})), sdk.TargetAny),
		b.MustMakeDef("path:to:ptrui64list", reflect.TypeOf(PtrU64List([]*uint64{})), sdk.TargetAny),
		// ptr float list
		b.MustMakeDef("path:to:ptrf32list", reflect.TypeOf(PtrF32List([]*float32{})), sdk.TargetAny),
		b.MustMakeDef("path:to:ptrf64list", reflect.TypeOf(PtrF64List([]*float64{})), sdk.TargetAny),
		// ptr complex list
		b.MustMakeDef("path:to:ptrc64list", reflect.TypeOf(PtrC64List([]*complex64{})), sdk.TargetAny),
		b.MustMakeDef("path:to:ptrc128list", reflect.TypeOf(PtrC128List([]*complex128{})), sdk.TargetAny),
		// ptr bool list
		b.MustMakeDef("path:to:ptrboollist", reflect.TypeOf(PtrBoolList([]*bool{})), sdk.TargetAny),
	}
	for _, def := range slices.Concat(defs, customDefs) {
		if err := reg.Define(def); err != nil {
			return nil, err
		}
	}
	return reg, nil
}
