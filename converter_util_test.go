package codemark

import (
	"reflect"
	"testing"
)

func TestTypeID(t *testing.T) {
	tests := []struct {
		name   string
		typ    any
		typeID string
	}{
		{
			name:   "pointer to slice",
			typ:    new(string),
			typeID: "ptr.string",
		},
		{
			name:   "slice of interfaces",
			typ:    []any{},
			typeID: "slice.interface",
		},
		{
			name:   "map",
			typ:    map[string]*int{},
			typeID: "map.string.ptr.int",
		},
		{
			name:   "chan",
			typ:    make(chan string),
			typeID: "chan.string",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rtype := reflect.TypeOf(tc.typ)
			typeID, err := TypeID(rtype)
			if err != nil {
				t.Errorf("err occured: %s\n", err)
			}
			if typeID != tc.typeID {
				t.Fatalf("expected: %s; got: %s\n", tc.typeID, typeID)
			}
		})
	}
}
