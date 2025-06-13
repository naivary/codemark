package testing

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
	"github.com/naivary/codemark/sdk/utils"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

type ValidValueFunc func(got, want reflect.Value) bool

type ConverterTesterConfig struct {
	ValidValueFuncs map[sdk.TypeID]ValidValueFunc
	Types           map[sdk.TypeID]reflect.Type
}

type ConverterTestCase struct {
	Name         string
	Marker       parser.Marker
	Target       sdk.Target
	ToType       reflect.Type
	IsValidCase  bool
	IsValidValue ValidValueFunc
}

type ConverterTester interface {
	Tests() ([]ConverterTestCase, error)
	Run(t *testing.T, tc ConverterTestCase, mngr sdk.ConverterManager)
}

type converterTester struct {
	conv  sdk.Converter
	vvfns map[sdk.TypeID]ValidValueFunc
	types map[sdk.TypeID]reflect.Type
}

func NewConverterTester(conv sdk.Converter, cfg *ConverterTesterConfig) (ConverterTester, error) {
	if cfg == nil {
		cfg = &ConverterTesterConfig{}
	}
	c := &converterTester{
		conv:  conv,
		vvfns: make(map[sdk.TypeID]ValidValueFunc),
		types: typeIDToReflectTypeMap(),
	}
	for typeID, fn := range cfg.ValidValueFuncs {
		vvfn := c.getVVFnFromTypeID(typeID)
		if vvfn != nil {
			return nil, fmt.Errorf("IsValidFunction exists: %s\n", typeID)
		}
		c.vvfns[typeID] = fn
	}
	for typeID, rtype := range cfg.Types {
		_, found := c.types[typeID]
		if found {
			return nil, fmt.Errorf("type id already exists: %s\n", typeID)
		}
		c.types[typeID] = rtype
	}
	return c, nil
}

func (c *converterTester) Tests() ([]ConverterTestCase, error) {
	types := c.conv.SupportedTypes()
	tests := make([]ConverterTestCase, 0, len(types))
	for _, rtype := range types {
		typeID := sdkutil.TypeIDOf(rtype)
		marker := RandMarker(rtype)
		if marker == nil {
			return nil, fmt.Errorf("no valid marker found: %v\n", rtype)
		}
		to := c.types[typeID]
		name := fmt.Sprintf("marker(%s) to %v", marker.Ident(), to)
		tc := ConverterTestCase{
			Name:         name,
			Marker:       marker,
			Target:       sdk.TargetAny,
			ToType:       to,
			IsValidCase:  true,
			IsValidValue: c.getVVFnFromTypeID(typeID),
		}
		tests = append(tests, tc)
	}
	return tests, nil
}

func (c *converterTester) Run(t *testing.T, tc ConverterTestCase, mngr sdk.ConverterManager) {
	t.Run(tc.Name, func(t *testing.T) {
		v, err := mngr.Convert(tc.Marker, tc.Target)
		if err != nil {
			t.Errorf("err occured: %s\n", err)
		}
		gotType := reflect.TypeOf(v)
		if gotType != tc.ToType {
			t.Fatalf("types don't match after conversion. got: %v; expected: %v\n", gotType, tc.ToType)
		}
		gotValue := reflect.ValueOf(v)
		if !tc.IsValidValue(gotValue, tc.Marker.Value()) {
			t.Fatalf("value is not correct. got: %v; wanted: %v\n", gotValue, tc.Marker.Value())
		}
		t.Logf("succesfully converted. got: %v; expected: %v\n", gotType, tc.ToType)
	})
}

func (c *converterTester) getVVFnFromTypeID(typeID string) ValidValueFunc {
	if sdkutil.MatchTypeID(typeID, `slice\..+`) {
		return c.isValidList
	}
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?int\d{0,2}`) {
		return isValidInteger
	}
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?float\d{2}`) {
		return isValidFloat
	}
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?complex\d{2,3}`) {
		return isValidComplex
	}
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?bool`) {
		return isValidBool
	}
	if sdkutil.MatchTypeID(typeID, `(ptr\.)?string`) {
		return isValidString
	}
	return c.vvfns[typeID]
}

func (c *converterTester) isValidList(got, want reflect.Value) bool {
	elem := got.Type().Elem()
	typeID := sdkutil.TypeIDOf(elem)
	vvfn := c.getVVFnFromTypeID(typeID)
	if vvfn == nil {
		return false
	}
	for i := 0; i < want.Len(); i++ {
		wantElem := want.Index(i)
		gotElem := got.Index(i)
		if !vvfn(gotElem, wantElem) {
			return false
		}
	}
	return true
}

func isValidInteger(got, want reflect.Value) bool {
	got = utils.DeRefValue(got)
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

func isValidFloat(got, want reflect.Value) bool {
	got = utils.DeRefValue(got)
	w := want.Interface().(float64)
	return AlmostEqual(got.Float(), w)
}

func isValidComplex(got, want reflect.Value) bool {
	got = utils.DeRefValue(got)
	w := want.Interface().(complex128)
	return got.Complex() == w
}

func isValidString(got, want reflect.Value) bool {
	got = utils.DeRefValue(got)
	w := want.Interface().(string)
	return got.String() == w
}

func isValidBool(got, want reflect.Value) bool {
	got = utils.DeRefValue(got)
	w := want.Interface().(bool)
	return got.Bool() == w
}

func typeIDToReflectTypeMap() map[sdk.TypeID]reflect.Type {
	types := DefaultTypes()
	m := make(map[sdk.TypeID]reflect.Type, len(types))
	for _, typ := range types {
		rtype := reflect.TypeOf(typ)
		typeID := sdkutil.TypeIDOf(rtype)
		m[typeID] = rtype
	}
	return m
}
