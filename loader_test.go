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
			name:  "simple file",
			paths: []string{"testdata/auth_req.go"},
		},
		{
			name:  "sub directory",
			paths: []string{"testdata/pkg1/consts.go"},
		},
	}

	reg := NewRegistry()
	defs := []*Definition{
		MakeDef("path:to:max", TargetConst, reflect.TypeOf(i(int(0)))),
		MakeDef("path:to:maxh", TargetInterface, reflect.TypeOf(i(int(0)))),
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
			infos, err := l.Load(tc.paths...)
			if err != nil {
				t.Fatal(err.Error())
			}

			for _, info := range infos {
				for _, iface := range info.Interfaces {
					t.Log(iface.Defs())
				}
			}
		})
	}
}
