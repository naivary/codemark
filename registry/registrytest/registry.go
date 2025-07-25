package registrytest

import (
	"reflect"
	"slices"

	optionv1 "github.com/naivary/codemark/api/option/v1"
	"github.com/naivary/codemark/marker/markertest"
	"github.com/naivary/codemark/option"
	"github.com/naivary/codemark/registry"
	"github.com/naivary/codemark/typeutil"
)

type (
	F32       float32
	F64       float64
	PtrF32    *float32
	PtrF64    *float64
	Bool      bool
	PtrBool   *bool
	C64       complex64
	C128      complex128
	PtrC64    *complex64
	PtrC128   *complex128
	String    string
	PtrString *string

	Int     int
	I8      int8
	I16     int16
	Byte    byte
	I32     int32
	I64     int64
	PtrInt  *int
	PtrI8   *int8
	PtrI16  *int16
	PtrByte *byte
	PtrI32  *int32
	PtrI64  *int64

	Uint    uint
	Rune    rune
	U8      uint8
	U16     uint16
	U32     uint32
	U64     uint64
	PtrUint *uint
	PtrRune *rune
	PtrU8   *uint8
	PtrU16  *uint16
	PtrU32  *uint32
	PtrU64  *uint64

	AnyList    []any
	StringList []string
	IntList    []int
	I8List     []int8
	I16List    []int16
	I32List    []int32
	ByteList   []byte
	I64List    []int64
	UintList   []uint
	RuneList   []rune
	U8List     []uint8
	U16List    []uint16
	U32List    []uint32
	U64List    []uint64
	F32List    []float32
	F64List    []float64
	C64List    []complex64
	C128List   []complex128
	BoolList   []bool

	PtrAnyList    []*any
	PtrStringList []*string
	PtrBoolList   []*bool
	PtrIntList    []*int
	PtrI8List     []*int8
	PtrI16List    []*int16
	PtrI32List    []*int32
	PtrByteList   []*byte
	PtrI64List    []*int64
	PtrUintList   []*uint
	PtrRuneList   []*rune
	PtrU8List     []*uint8
	PtrU16List    []*uint16
	PtrU32List    []*uint32
	PtrU64List    []*uint64
	PtrF32List    []*float32
	PtrF64List    []*float64
	PtrC64List    []*complex64
	PtrC128List   []*complex128
)

func FloatTypes() []any {
	return []any{
		F32(0), F64(0),
		PtrF32(nil), PtrF64(nil),
	}
}

func ComplexTypes() []any {
	return []any{
		C64(0 + 0i), C128(0 + 0i),
		PtrC64(nil), PtrC128(nil),
	}
}

func BoolTypes() []any {
	return []any{
		Bool(false),
		PtrBool(nil),
	}
}

func StringTypes() []any {
	return []any{
		String(""),
		PtrString(nil),
	}
}

func IntTypes() []any {
	return []any{
		Int(0), I8(0), I16(0), I32(0), I64(0),
		PtrInt(nil), PtrI8(nil), PtrI16(nil), PtrI32(nil), PtrI64(nil),
	}
}

func UintTypes() []any {
	return []any{
		Uint(0), U8(0), U16(0), U32(0), U64(0),
		PtrUint(nil), PtrU8(nil), PtrU16(nil), PtrU32(nil), PtrU64(nil),
	}
}

func ListTypes() []any {
	return []any{
		StringList(nil), IntList(nil), I8List(nil), I16List(nil), I32List(nil), I64List(nil),
		UintList(nil), U8List(nil), U16List(nil), U32List(nil), U64List(nil),
		F32List(nil), F64List(nil), C64List(nil), C128List(nil), BoolList(nil), AnyList(nil),
		PtrStringList(nil), PtrBoolList(nil),
		PtrIntList(nil), PtrI8List(nil), PtrI16List(nil), PtrI32List(nil), PtrI64List(nil),
		PtrUintList(nil), PtrU8List(nil), PtrU16List(nil), PtrU32List(nil), PtrU64List(nil),
		PtrF32List(nil), PtrF64List(nil), PtrC64List(nil), PtrC128List(nil), PtrAnyList(nil),
	}
}

func AllTypes() []any {
	return slices.Concat(
		FloatTypes(),
		ComplexTypes(),
		BoolTypes(),
		StringTypes(),
		IntTypes(),
		UintTypes(),
		ListTypes(),
	)
}

// AliasOpts are go-native types which are supported by the default
// converters but cannot be registered using the native types because they are
// aliases of others e.g. byte=uint8 and rune=int32.
func AliasOpts() []optionv1.Option {
	return []optionv1.Option{
		option.MustMake(
			markertest.NewIdent("byte"),
			reflect.TypeOf(Byte(0)),
			nil, false,
			optionv1.TargetAny,
		),
		option.MustMake(
			markertest.NewIdent("rune"),
			reflect.TypeOf(Rune(0)),
			nil, false,
			optionv1.TargetAny,
		),
		option.MustMake(
			markertest.NewIdent("ptr.byte"),
			reflect.TypeOf(PtrByte(nil)),
			nil, false,
			optionv1.TargetAny,
		),
		option.MustMake(
			markertest.NewIdent("ptr.rune"),
			reflect.TypeOf(PtrRune(nil)),
			nil, false,
			optionv1.TargetAny,
		),
		option.MustMake(
			markertest.NewIdent("slice.byte"),
			reflect.TypeOf(ByteList(nil)),
			nil, false,
			optionv1.TargetAny,
		),
		option.MustMake(
			markertest.NewIdent("slice.rune"),
			reflect.TypeOf(RuneList(nil)),
			nil, false,
			optionv1.TargetAny,
		),
		option.MustMake(
			markertest.NewIdent("slice.ptr.byte"),
			reflect.TypeOf(PtrByteList(nil)),
			nil, false,
			optionv1.TargetAny,
		),
		option.MustMake(
			markertest.NewIdent("slice.ptr.rune"),
			reflect.TypeOf(PtrRuneList(nil)),
			nil, false,
			optionv1.TargetAny,
		),
	}
}

// NewOptsSet returns all the options matching the default markers available with
// the correct type
func NewOptsSet() []optionv1.Option {
	types := AllTypes()
	aliases := AliasOpts()
	opts := make([]optionv1.Option, 0, len(types)+len(aliases))
	for _, typ := range AllTypes() {
		rtype := reflect.TypeOf(typ)
		name := typeutil.NameFor(rtype)
		ident := markertest.NewIdent(name)
		opt := option.MustMake(ident, rtype, nil, false, optionv1.TargetAny)
		opts = append(opts, opt)
	}
	return slices.Concat(opts, aliases)
}

func NewRegistry(opts []optionv1.Option) (registry.Registry, error) {
	reg := registry.InMemory()
	for _, opt := range opts {
		if err := reg.Define(&opt); err != nil {
			return nil, err
		}
	}
	return reg, nil
}
