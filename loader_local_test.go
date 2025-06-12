package codemark

import (
	"reflect"
	"testing"

	"github.com/naivary/codemark/sdk"
)

type description string

type required bool

func localLoaderDefs() sdk.Registry {
	reg := NewInMemoryRegistry()
	defs := []*sdk.Definition{
		MustMakeDef("openapi_v3:general:description", reflect.TypeOf(description("")), sdk.TargetStruct),
		MustMakeDef("openapi_v3:validation:required", reflect.TypeOf(required(false)), sdk.TargetField),
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
			_ = projs
		})
	}
}
