package testing

import (
	"reflect"
	"slices"

	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
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

	StringList    []string
	IntList       []int
	I8List        []int8
	I16List       []int16
	I32List       []int32
	ByteList      []byte
	I64List       []int64
	UintList      []uint
	RuneList      []rune
	U8List        []uint8
	U16List       []uint16
	U32List       []uint32
	U64List       []uint64
	F32List       []float32
	F64List       []float64
	C64List       []complex64
	C128List      []complex128
	BoolList      []bool
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
		F32List(nil), F64List(nil), C64List(nil), C128List(nil), BoolList(nil),
		PtrStringList(nil), PtrBoolList(nil),
		PtrIntList(nil), PtrI8List(nil), PtrI16List(nil), PtrI32List(nil), PtrI64List(nil),
		PtrUintList(nil), PtrU8List(nil), PtrU16List(nil), PtrU32List(nil), PtrU64List(nil),
		PtrF32List(nil), PtrF64List(nil), PtrC64List(nil), PtrC128List(nil),
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

// NewRegistry returns a registry which includes definitions for all the
// default types which can be converterd using the builtin converter.
func NewRegistry(reg sdk.Registry, b sdk.DefinitionMaker, customDefs ...*sdk.Definition) (sdk.Registry, error) {
	for _, typ := range AllTypes() {
		rtype := reflect.TypeOf(typ)
		name := sdkutil.NameFor(rtype)
		ident := NewIdent(name)
		def, err := b.MakeDef(ident, rtype, sdk.TargetAny)
		if err != nil {
			return nil, err
		}
		if err := reg.Define(def); err != nil {
			return nil, err
		}
	}
	// add special definitions which cannot be added using NameFor because they
	// are aliases to other types e.g. byte=int32 and rune=uint8
	defs := []*sdk.Definition{
		b.MustMakeDef(NewIdent("byte"), reflect.TypeOf(Byte(0)), sdk.TargetAny),
		b.MustMakeDef(NewIdent("rune"), reflect.TypeOf(Rune(0)), sdk.TargetAny),
		b.MustMakeDef(NewIdent("ptr.byte"), reflect.TypeOf(PtrByte(nil)), sdk.TargetAny),
		b.MustMakeDef(NewIdent("ptr.rune"), reflect.TypeOf(PtrRune(nil)), sdk.TargetAny),
		b.MustMakeDef(NewIdent("slice.byte"), reflect.TypeOf(ByteList(nil)), sdk.TargetAny),
		b.MustMakeDef(NewIdent("slice.rune"), reflect.TypeOf(RuneList(nil)), sdk.TargetAny),
		b.MustMakeDef(NewIdent("slice.ptr.byte"), reflect.TypeOf(PtrByteList(nil)), sdk.TargetAny),
		b.MustMakeDef(NewIdent("slice.ptr.rune"), reflect.TypeOf(PtrRuneList(nil)), sdk.TargetAny),
	}
	for _, def := range slices.Concat(defs, customDefs) {
		if err := reg.Define(def); err != nil {
			return nil, err
		}
	}
	return reg, nil
}
