package codemark

import (
	"reflect"
	"testing"
)

var strconv = &stringConverter{}

func TestNewConvMngr(t *testing.T) {
	tests := []struct {
		name  string
		conv  Converter
		types []reflect.Type
	}{
		{
			name:  "adding converter",
			conv:  strconv,
			types: strconv.SupportedTypes(),
		},
	}

	reg := NewInMemoryRegistry()
	def := MakeDef("path:to:marker", TargetField, reflect.TypeOf(bool(false)))
	if err := reg.Define(def); err != nil {
		t.Errorf("err occured: %s\n", err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mngr, err := NewConvMngr(reg, tc.conv)
			if err != nil {
				t.Errorf("err occured: %s\n", err)
			}
			for _, rtype := range tc.types {
				typeID, err := TypeID(rtype)
				if err != nil {
					t.Errorf("err occured: %s\n", err)
				}
				conv, err := mngr.GetConverter(rtype)
				if err != nil {
					t.Errorf("err occured: %s\n", err)
				}
				if conv == nil || conv != strconv {
					t.Fatalf("converter for type ID `%s` must exist but is not", typeID)
				}
			}
		})
	}
}
