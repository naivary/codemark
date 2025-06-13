package codemark

import (
	"reflect"
	"testing"

	sdktesting "github.com/naivary/codemark/sdk/testing"
)

func TestListConverter(t *testing.T) {
	convTester, err := sdktesting.NewConverterTester(nil, nil)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	tests, err := convTester.NewTest(&listConverter{})
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
			if err != nil && tc.IsValidCase {
				t.Errorf("err occured: %s\n", err)
			}
			if err != nil && !tc.IsValidCase {
				t.Skipf("wanted err occured: %v\n", err)
			}
			rtype := reflect.TypeOf(v)
			if rtype != tc.ToType {
				t.Fatalf("types don't match after conversion. got: %v; expected: %v\n", rtype, tc.ToType)
			}
			rvalue := reflect.ValueOf(v)
			if !tc.IsValidValue(rvalue, tc.Marker.Value()) {
				t.Fatalf("value is not correct. got: %v; wanted: %v\n", rvalue, tc.Marker.Value())
			}
			t.Logf("succesfully converted. got: %v; expected: %v\n", rtype, tc.ToType)
		})
	}
}
