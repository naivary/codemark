package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

const _boolValue = true

func TestBoolConverter(t *testing.T) {
	tests := []struct {
		name         string
		mrk          parser.Marker
		t            sdk.Target
		out          reflect.Type
		isValid      bool
		isValidValue func(got reflect.Value) bool
	}{
		{
			name:    "bool marker to bool type",
			mrk:     parser.NewMarker("path:to:bool", parser.MarkerKindBool, reflect.ValueOf(_boolValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.Bool(false)),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				b := got.Interface().(sdktesting.Bool)
				return b == _boolValue
			},
		},
		{
			name:    "bool marker to ptr bool type",
			mrk:     parser.NewMarker("path:to:ptrbool", parser.MarkerKindBool, reflect.ValueOf(_boolValue)),
			t:       sdk.TargetField,
			out:     reflect.TypeOf(sdktesting.PtrBool(new(bool))),
			isValid: true,
			isValidValue: func(got reflect.Value) bool {
				b := got.Interface().(sdktesting.PtrBool)
				return *b == _boolValue
			},
		},
	}
	reg, err := sdktesting.NewDefsSet(NewInMemoryRegistry(), &DefinitionMarker{})
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	mngr, err := NewConvMngr(reg)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, err := mngr.Convert(tc.mrk, tc.t)
			if err != nil {
				t.Errorf("err occured: %s\n", err)
			}
			rtype := reflect.TypeOf(v)
			if rtype != tc.out {
				t.Fatalf("types don't match after conversion. got: %v; expected: %v\n", rtype, tc.out)
			}
			rvalue := reflect.ValueOf(v)
			if !tc.isValidValue(rvalue) {
				t.Fatalf("value is not correct. got: %v", rvalue)
			}
			t.Logf("succesfully converted. got: %v; expected: %v\n", rtype, tc.out)
		})
	}
}
