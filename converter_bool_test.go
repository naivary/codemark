package codemark

import (
	"reflect"
	"testing"

	sdktesting "github.com/naivary/codemark/sdk/testing"
)

func TestBoolConverter(t *testing.T) {
	convTester, err := sdktesting.NewConverterTester(nil, nil)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	tests, err := convTester.NewTest(&boolConverter{})
	if err != nil {
		t.Errorf("err occured: %s\n", err)
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
