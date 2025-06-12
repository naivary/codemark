package testing

import (
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
		typ     any
		isValid func(got any) bool
	}{
		{
			name:    "string marker",
			typ:     string(""),
			isValid: isNotEmpty[string],
		},
		{
			name:    "int marker",
			typ:     int(0),
			isValid: isValidNumber[int64],
		},
		{
			name:    "int16 marker",
			typ:     int16(0),
			isValid: isValidNumber[int64],
		},
		{
			name:    "float32 marker",
			typ:     float32(0.0),
			isValid: isValidNumber[float64],
		},
		{
			name:    "float64 marker",
			typ:     float64(0.0),
			isValid: isValidNumber[float64],
		},

		{
			name:    "ptr string marker",
			typ:     new(string),
			isValid: isNotEmpty[string],
		},
		{
			name:    "list string marker",
			typ:     []string{},
			isValid: isNotEmpty[[]any],
		},
		{
			name:    "list int marker",
			typ:     []int{},
			isValid: isNotEmpty[[]any],
		},
		{
			name:    "list float marker",
			typ:     []float64{},
			isValid: isNotEmpty[[]any],
		},
		{
			name:    "list ptr float marker",
			typ:     new([]float64),
			isValid: isNotEmpty[[]any],
		},

		{
			name:    "unknown marker",
			typ:     new([]string),
			isValid: isNotEmpty[[]any],
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			marker := RandMarker(tc.typ)
			if marker == nil {
				t.Errorf("err: no marker produced")
			}
			val := marker.Value().Interface()
			if !tc.isValid(val) {
				t.Errorf("value is not valid. got: %v\n", val)
			}
		})
	}
}
