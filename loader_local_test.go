package codemark

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/naivary/codemark/sdk"
)

type description string

type required bool

func localLoaderDefs() sdk.Registry {
	reg := NewInMemoryRegistry()
	defs := []*sdk.Definition{
		MustMakeDef("openapi_v3:general:description", sdk.TargetStruct, reflect.TypeOf(description(""))),
		MustMakeDef("openapi_v3:validation:required", sdk.TargetField, reflect.TypeOf(required(false))),
	}
	for _, def := range defs {
		if err := reg.Define(def); err != nil {
			panic(err)
		}
	}
	return reg
}

func TestLocalLoader(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
	}{
		{
			name:    "default behavior",
			pattern: "./testdata",
		},
	}
	reg := localLoaderDefs()
	mngr, err := NewConvMngr(reg)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	l := NewLocalLoader(mngr, nil)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			projs, err := l.Load(tc.pattern)
			if err != nil {
				t.Errorf("err occured: %s\n", err)
			}
			for _, proj := range projs {
				for _, stru := range proj.Structs {
					for _, def := range stru.Defs {
						fmt.Println(def[0].(description))
					}
				}
			}
		})
	}
}
