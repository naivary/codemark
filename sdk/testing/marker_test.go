package testing

import (
	"reflect"
	"testing"

	"golang.org/x/exp/constraints"
)

type lengthy interface {
	~string | []any
}

func isNotEmpty[T lengthy](v any) bool {
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
		ident   string
		isValid func(got any) bool
	}{
		{
			name:    "string marker",
			rtype:   reflect.TypeOf(string("")),
			ident:   "path:to:str",
			isValid: isNotEmpty[string],
		},
		{
			name:    "int marker",
			rtype:   reflect.TypeOf(int(0)),
			ident:   "int",
			isValid: isValidNumber[int64],
		},
		{
			name:    "float32 marker",
			rtype:   reflect.TypeOf(float32(0.0)),
			ident:   "float32",
			isValid: isValidNumber[float64],
		},
		{
			name:    "float64 marker",
			rtype:   reflect.TypeOf(float64(0.0)),
			ident:   "float64",
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
			name:    "unknown marker",
			rtype:   reflect.TypeOf(new([]string)),
			isValid: isNotEmpty[[]any],
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			marker := RandomMarker(tc.rtype)
			ident := marker.Ident()
			if ident != tc.ident {
				t.Errorf("identifier dont macht. got: %s; wanted: %s\n", ident, tc.ident)
			}
			val := marker.Value().Interface()
			if !tc.isValid(val) {
				t.Errorf("value is not valid. got: %v\n", val)
			}
		})
	}
}
