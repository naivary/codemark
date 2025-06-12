package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/parser"
	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

func TestStringConverter(t *testing.T) {
	tests := []struct {
		name          string
		mrk           parser.Marker
		t             sdk.Target
		out           reflect.Type
		expectedValue any
		isValid       bool
		isValidValue  func(got reflect.Value) bool
	}{
		{
			name:          "string marker to string type",
			mrk:           parser.NewMarker("path:to:str", parser.MarkerKindString, reflect.ValueOf(string("codemark"))),
			t:             sdk.TargetField,
			out:           reflect.TypeOf(sdktesting.String("")),
			isValid:       true,
			expectedValue: "codemark",
			isValidValue: func(got reflect.Value) bool {
				str := got.Interface().(sdktesting.String)
				return str == "codemark"
			},
		},
		{
			name:          "string marker to ptr string type",
			mrk:           parser.NewMarker("path:to:ptrstr", parser.MarkerKindString, reflect.ValueOf(string("codemark"))),
			t:             sdk.TargetField,
			out:           reflect.TypeOf(sdktesting.PtrString(new(string))),
			expectedValue: "codemark",
			isValid:       true,
			isValidValue: func(got reflect.Value) bool {
				str := got.Interface().(sdktesting.PtrString)
				return *str == "codemark"
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
