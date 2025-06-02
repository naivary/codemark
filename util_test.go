package codemark

import (
	"reflect"
	"strings"
	"testing"
)

func TestKindPath(t *testing.T) {
	tests := []struct {
		name     string
		input    reflect.Type
		expected string
	}{
		{
			name:     "slice",
			input:    reflect.TypeOf([]string{}),
			expected: "slice.string",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			id, err := TypeID(tc.input, &b)
			if err != nil {
			}
			if id != tc.expected {
				t.Fatalf("expected: %s. got: %s", tc.expected, id)
			}
		})
	}
}
