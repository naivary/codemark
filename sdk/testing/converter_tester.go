package testing

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/naivary/codemark/sdk"
	sdkutil "github.com/naivary/codemark/sdk/utils"
)

type converterTester struct {
	conv  sdk.Converter
	vvfns map[sdk.TypeID]ValidValueFunc
	types map[sdk.TypeID]reflect.Type
}

func NewConverterTester(conv sdk.Converter) (ConverterTester, error) {
	c := &converterTester{
		conv:  conv,
		vvfns: make(map[sdk.TypeID]ValidValueFunc),
		types: make(map[sdk.TypeID]reflect.Type),
	}
	return c, nil
}

func (c *converterTester) ValidTests() ([]ConverterTestCase, error) {
	types := c.conv.SupportedTypes()
	tests := make([]ConverterTestCase, 0, len(types))
	for _, rtype := range types {
		tc, err := c.NewTest(rtype, true, sdk.TargetAny)
		if err != nil {
			return nil, err
		}
		tests = append(tests, tc)
	}
	return tests, nil
}

func (c *converterTester) NewTest(rtype reflect.Type, isValidCase bool, target sdk.Target) (ConverterTestCase, error) {
	typeID := sdkutil.TypeIDOf(rtype)
	marker := RandMarker(rtype)
	if marker == nil {
		return ConverterTestCase{}, fmt.Errorf("no valid marker found: %v\n", rtype)
	}
	to := c.types[typeID]
	name := fmt.Sprintf("marker[%s] to %v", marker.Ident(), to)
	tc := ConverterTestCase{
		Name:         name,
		Marker:       marker,
		Target:       target,
		ToType:       to,
		IsValidCase:  isValidCase,
		IsValidValue: c.vvfns[typeID],
	}
	return tc, nil
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

func (c *converterTester) AddVVFunc(rtype reflect.Type, fn ValidValueFunc) error {
	typeID := sdkutil.TypeIDOf(rtype)
	_, found := c.vvfns[typeID]
	if found {
		return fmt.Errorf("ValidValueFunc already exists: %s\n", typeID)
	}
	c.vvfns[typeID] = fn
	return nil
}

func (c *converterTester) AddType(rtype reflect.Type) error {
	typeID := sdkutil.TypeIDOf(rtype)
	_, found := c.types[typeID]
	if found {
		return fmt.Errorf("reflect.Type already exists: %s\n", typeID)
	}
	c.types[typeID] = rtype
	return nil
}
