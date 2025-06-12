package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/sdk"
	sdktesting "github.com/naivary/codemark/sdk/testing"
)

const _complexValue = 9 + 9i

func TestComplexConverter(t *testing.T) {
	tests := []sdktesting.ConverterTestCase{
		{
			Name:        "complex marker to c64 type",
			Marker:      sdktesting.RandMarker(sdktesting.C64(0 + 0i)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.C64(0 + 0i)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := complex128(got.Interface().(sdktesting.C64))
				w := wanted.Interface().(complex128)
				return g == w
			},
		},
		{
			Name:        "complex marker to c128 type",
			Marker:      sdktesting.RandMarker(sdktesting.C128(0 + 0i)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.C128(0 + 0i)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := complex128(got.Interface().(sdktesting.C128))
				w := wanted.Interface().(complex128)
				return g == w
			},
		},
		{
			Name:        "complex marker to ptrc64 type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrC64(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrC64(nil)),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := complex128(*got.Interface().(sdktesting.PtrC64))
				w := wanted.Interface().(complex128)
				return g == w
			},
		},
		{
			Name:        "complex marker to ptrc128 type",
			Marker:      sdktesting.RandMarker(sdktesting.PtrC128(nil)),
			Target:      sdk.TargetAny,
			ToType:      reflect.TypeOf(sdktesting.PtrC128(new(complex128))),
			IsValidCase: true,
			IsValidValue: func(got, wanted reflect.Value) bool {
				g := complex128(*got.Interface().(sdktesting.PtrC128))
				w := wanted.Interface().(complex128)
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
