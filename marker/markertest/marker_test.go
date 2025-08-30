package markertest

import (
	"reflect"
	"testing"

	"golang.org/x/exp/constraints"
)

func TestRandMarker(t *testing.T) {
	tests := []struct {
		name        string
		typ         any
		isValid     func(got any) bool
		isValidCase bool
	}{
		{
			name:        "string marker",
			typ:         string(""),
			isValid:     isNotEmpty[string],
			isValidCase: true,
		},
		{
			name:        "int marker",
			typ:         int(0),
			isValid:     isValidNumber[int64],
			isValidCase: true,
		},
		{
			name:        "int16 marker",
			typ:         int16(0),
			isValid:     isValidNumber[int64],
			isValidCase: true,
		},
		{
			name:        "float32 marker",
			typ:         float32(0.0),
			isValid:     isValidNumber[float64],
			isValidCase: true,
		},
		{
			name:        "float64 marker",
			typ:         float64(0.0),
			isValid:     isValidNumber[float64],
			isValidCase: true,
		},

		{
			name:        "ptr string marker",
			typ:         new(string),
			isValid:     isNotEmpty[string],
			isValidCase: true,
		},
		{
			name:        "list string marker",
			typ:         []string{},
			isValid:     isNotEmpty[[]any],
			isValidCase: true,
		},
		{
			name:    "list int marker",
			typ:     []int{},
			isValid: isNotEmpty[[]any],

			isValidCase: true,
		},
		{
			name:        "list float marker",
			typ:         []float64{},
			isValid:     isNotEmpty[[]any],
			isValidCase: true,
		},
		{
			name:        "list ptr float marker",
			typ:         []*float64{},
			isValid:     isNotEmpty[[]any],
			isValidCase: true,
		},
		{
			name:        "unknown marker",
			typ:         new([]string),
			isValidCase: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rtype := reflect.TypeOf(tc.typ)
			marker, err := Rand(rtype)
			if marker == nil && tc.isValidCase {
				t.Errorf("err: no marker produced")
			}
			if err == nil && !tc.isValidCase {
				t.Fatalf("err: expected to fail but didnt: %v\n", marker)
			}
			if err != nil && !tc.isValidCase {
				t.Skipf("valid test case was recognised correctly. Skipping the rest")
			}
			val := marker.Value.Interface()
			if !tc.isValid(val) {
				t.Errorf("value is not valid. got: %v\n", val)
			}
		})
	}
}

func isNotEmpty[T ~string | []any](v any) bool {
	return len(v.(T)) != 0
}

func isValidNumber[T constraints.Integer | constraints.Float](got any) bool {
	i := got.(T)
	return i >= 0
}
