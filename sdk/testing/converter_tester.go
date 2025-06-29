package testing

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/naivary/codemark/definition/target"
	"github.com/naivary/codemark/sdk"
)

var _ ConverterTester = (*converterTester)(nil)

type converterTester struct {
	conv    sdk.Converter
	vvfns   map[reflect.Type]ValidValueFunc
	toTypes map[reflect.Type]reflect.Type
}

// NewConvTester returns a new ConverterTester for the given converter. The
// parameter `toTypes` is providing a map of a supported type to an example
// custom type which is being used by the converter as a test. For example the
// built in integer converter is converter an int to a type Int int.
func NewConvTester(conv sdk.Converter, toTypes map[reflect.Type]reflect.Type) (ConverterTester, error) {
	if toTypes == nil {
		return nil, errors.New("missing toTypes map")
	}
	c := &converterTester{
		conv:    conv,
		vvfns:   make(map[reflect.Type]ValidValueFunc),
		toTypes: toTypes,
	}
	return c, nil
}

func (c *converterTester) NewTest(from reflect.Type, isValidCase bool, t target.Target) (ConverterTestCase, error) {
	marker, err := RandMarker(from)
	if err != nil {
		return ConverterTestCase{}, err
	}
	to := c.toTypes[from]
	if to == nil {
		return ConverterTestCase{}, fmt.Errorf("no to reflect.Type found: %v\n", from)
	}
	name := fmt.Sprintf("marker[%s] to %v", marker.Ident(), to)
	vvfn := c.vvfns[to]
	if vvfn == nil {
		return ConverterTestCase{}, fmt.Errorf("ValidValueFunc not found: %v\n", from)
	}
	tc := ConverterTestCase{
		Name:         name,
		Marker:       *marker,
		Target:       t,
		To:           to,
		IsValidCase:  isValidCase,
		IsValidValue: vvfn,
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
		if gotType != tc.To {
			t.Fatalf("types don't match after conversion. got: %v; expected: %v\n", gotType, tc.To)
		}
		gotValue := reflect.ValueOf(v)
		if !tc.IsValidValue(gotValue, tc.Marker.Value()) {
			t.Fatalf("value is not correct. got: %v; wanted: %v\n", gotValue, tc.Marker.Value())
		}
		t.Logf("succesfully converted. got: %v; expected: %v\n", gotType, tc.To)
	})
}

func (c *converterTester) AddVVFunc(to reflect.Type, fn ValidValueFunc) error {
	_, found := c.vvfns[to]
	if found {
		return fmt.Errorf("ValidValueFunc already exists: %s\n", to)
	}
	c.vvfns[to] = fn
	return nil
}

func (c *converterTester) ValidTests() ([]ConverterTestCase, error) {
	types := c.conv.SupportedTypes()
	tests := make([]ConverterTestCase, 0, len(types))
	for _, from := range types {
		tc, err := c.NewTest(from, true, target.ANY)
		if err != nil {
			return nil, err
		}
		tests = append(tests, tc)
	}
	return tests, nil
}
