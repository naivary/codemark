package typeutil

import (
	"reflect"
	"testing"
)

func TestNameFor(t *testing.T) {
	tests := []struct {
		name     string
		typ      any
		expected string
	}{
		{
			name:     "pointer to slice",
			typ:      new(string),
			expected: "ptr.string",
		},
		{
			name:     "slice of interfaces",
			typ:      []any{},
			expected: "slice.interface",
		},
		{
			name:     "map",
			typ:      map[string]*int{},
			expected: "map.string.ptr.int",
		},
		{
			name:     "chan",
			typ:      make(chan string),
			expected: "chan.string",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rtype := reflect.TypeOf(tc.typ)
			name := NameFor(rtype)
			if name != tc.expected {
				t.Fatalf("expected: %s; got: %s\n", tc.expected, name)
			}
		})
	}
}
