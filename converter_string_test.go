package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

func TestStringConverter(t *testing.T) {
	tests := []sdktesting.ConverterTestCase{
		{
			Name:        "string marker to string type",
			Marker:      sdktesting.RandMarker(sdktesting.String("")),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.String("")),
			IsValidCase: true,
			IsValidValue: func(got, want reflect.Value) bool {
				g := string(got.Interface().(sdktesting.String))
				w := want.Interface().(string)
				return g == w
			},
		},
		{
			Name:        "string marker to ptr string type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrString(nil)),
			Target:      sdk.TargetField,
			ToType:      reflect.TypeOf(sdktesting.PtrString(nil)),
			IsValidCase: true,
			IsValidValue: func(got, want reflect.Value) bool {
				g := string(*got.Interface().(sdktesting.PtrString))
				w := want.Interface().(string)
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
				t.Fatalf("value is not correct. got: %v; wanted: %v\n", gotValue, tc.Marker.Value())
			}
			t.Logf("succesfully converted. got: %v; expected: %v\n", gotType, tc.ToType)
		})
	}
}
