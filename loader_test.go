package main

import (
	"reflect"
	"testing"
)

func TestLoader(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
	}{
		{
			name:  "single file",
			paths: []string{"testdata/auth_req.go"},
		},
		{
			name:  "recursive",
			paths: []string{"./..."},
		},
	}

	reg := NewRegistry()
	defs := []*Definition{
		MakeDef("path:to:const", TargetConst, reflect.TypeOf(i(int(0)))),
		MakeDef("path:to:iface", TargetInterface, reflect.TypeOf(i(int(0)))),
		MakeDef("path:to:func", TargetFunc, reflect.TypeOf(i(int(0)))),
		MakeDef("path:to:field", TargetField, reflect.TypeOf(i(int(0)))),
		MakeDef("path:to:pkg", TargetPackage, reflect.TypeOf(str(""))),
	}
	for _, def := range defs {
		if err := reg.Define(def); err != nil {
			t.Error(err)
		}
	}
	for _, tc := range tests {
		conv, err := NewConverter(reg)
		if err != nil {
			t.Error(err)
		}
		t.Run(tc.name, func(t *testing.T) {
			l := NewLoader(conv, nil)
			_, err := l.Load(tc.paths...)
			if err != nil {
				t.Fatal(err.Error())
			}
		})
	}
}
