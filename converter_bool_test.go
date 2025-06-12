package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

func TestBoolConverter(t *testing.T) {
	tests := []sdktesting.ConverterTestCase{
		{
			Name:        "bool marker to bool type",
			Marker:      sdktesting.RandMarker(sdktesting.Bool(false)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.Bool(false)),
			IsValidCase: true,
			IsValidValue: func(got, want reflect.Value) bool {
				g := bool(got.Interface().(sdktesting.Bool))
				w := want.Interface().(bool)
				return w == g
			},
		},
		{
			Name:        "bool marker to ptr bool type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrBool(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrBool(nil)),
			IsValidCase: true,
			IsValidValue: func(got, want reflect.Value) bool {
				g := bool(*got.Interface().(sdktesting.PtrBool))
				w := want.Interface().(bool)
				return g == w
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
				t.Fatalf("value is not correct. got: %v; wanted: %v", gotValue, tc.Marker.Value())
			}
			t.Logf("succesfully converted. got: %v; expected: %v\n", gotType, tc.ToType)
		})
	}
}
