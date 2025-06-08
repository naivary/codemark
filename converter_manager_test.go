package codemark

import (
	"reflect"
	"testing"
)

var strConv = &stringConverter{}

func TestNewConvMngr(t *testing.T) {
	tests := []struct {
		name  string
		conv  Converter
		types []reflect.Type
	}{
		{
			name:  "adding converter",
			conv:  strConv,
			types: strConv.SupportedTypes(),
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
				typeID := TypeID(rtype)
				conv, err := mngr.GetConverter(rtype)
				if err != nil {
					t.Errorf("err occured: %s\n", err)
				}
				if conv == nil || conv != strConv {
					t.Fatalf("converter for type ID `%s` must exist but is not", typeID)
				}
			}
		})
	}
}
