package testing

import (
	"reflect"
	"testing"

	"golang.org/x/exp/constraints"
)

func isNotEmpty[T ~string | []any](v any) bool {
	return len(v.(T)) != 0
}

func isValidNumber[T constraints.Integer | constraints.Float](got any) bool {
	i := got.(T)
	return i >= 0
}

func TestRandMarker(t *testing.T) {
	tests := []struct {
		name    string
		rtype   reflect.Type
		isValid func(got any) bool
	}{
		{
			name:    "string marker",
			rtype:   reflect.TypeOf(string("")),
			isValid: isNotEmpty[string],
		},
		{
			name:    "int marker",
			rtype:   reflect.TypeOf(int(0)),
			isValid: isValidNumber[int64],
		},
		{
			name:    "float32 marker",
			rtype:   reflect.TypeOf(float32(0.0)),
			isValid: isValidNumber[float64],
		},
		{
			name:    "float64 marker",
			rtype:   reflect.TypeOf(float64(0.0)),
			isValid: isValidNumber[float64],
		},

		{
			name:    "ptr string marker",
			rtype:   reflect.TypeOf(new(string)),
			isValid: isNotEmpty[string],
		},
		{
			name:    "list string marker",
			rtype:   reflect.TypeOf([]string{}),
			isValid: isNotEmpty[[]any],
		},
		{
			name:    "list int marker",
			rtype:   reflect.TypeOf([]int{}),
			isValid: isNotEmpty[[]any],
		},
		{
			name:    "list float marker",
			rtype:   reflect.TypeOf([]float64{}),
			isValid: isNotEmpty[[]any],
		},
		{
			name:    "list ptr float marker",
			rtype:   reflect.TypeOf(new([]float64)),
			isValid: isNotEmpty[[]any],
		},

		{
			name:    "unknown marker",
			rtype:   reflect.TypeOf(new([]string)),
			isValid: isNotEmpty[[]any],
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			marker := RandMarker(tc.rtype)
			val := marker.Value().Interface()
			t.Log(val)
			if !tc.isValid(val) {
				t.Errorf("value is not valid. got: %v\n", val)
			}
		})
	}
}
