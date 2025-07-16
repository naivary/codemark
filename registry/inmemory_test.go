package registry

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/naivary/codemark/api"
	"github.com/naivary/codemark/definition/target"
	"github.com/naivary/codemark/maker"
)

type registryTestCase struct {
	// Name of the test case
	Name string
	// The definition being tested
	Def *api.Definition
	// Whether the test case is checking correctness or not.
	IsValid bool
}

type registryTester struct {
	reg Registry
}

func newRegTester(reg Registry) *registryTester {
	return &registryTester{
		reg: reg,
	}
}

func (r *registryTester) newTest(def *api.Definition, isValid bool) registryTestCase {
	return registryTestCase{
		Name:    fmt.Sprintf("%s[%s]", def.Ident, def.Output),
		IsValid: isValid,
		Def:     def,
	}
}

func (r *registryTester) run(t *testing.T, tc registryTestCase) {
	err := r.reg.Define(tc.Def)
	if err == nil && !tc.IsValid {
		t.Errorf("expected an error but err was nil: %s\n", tc.Def.Ident)
	}
	if err != nil && !tc.IsValid {
		t.Skipf("expected error to occur because it's an invalid test case: %s. Skipping...", err)
	}
	if err != nil && tc.IsValid {
		t.Errorf("could not define definition: %v\n", tc.Def)
	}
	def, err := r.reg.Get(tc.Def.Ident)
	if err != nil {
		t.Errorf("get failed with an error: %s\n", err)
	}
	if def != tc.Def {
		t.Errorf("definitions are not equal after retrieval. got: %v\n want: %v\n", def, tc.Def)
	}
	r.validateDoc(t, tc.Def, def)
	t.Logf("test case sucessfull: %s\n", tc.Name)
}

func (r *registryTester) validateDoc(t *testing.T, got, want *api.Definition) {
	if got.Doc == nil {
		t.Logf("no assertions will be done for the documentation because `%s` has no doc", got.Ident)
		return
	}
	doc, err := r.reg.DocOf(got.Ident)
	if err != nil {
		t.Errorf("DocOf failed with an error: %s\n", err)
	}
	if doc != want.Doc {
		t.Errorf("doc is not the same. got: %s; want: %s\n", doc, want.Doc)
	}
}

func defs() []*api.Definition {
	return []*api.Definition{
		maker.MustMakeDef("codemark:registry:plain", reflect.TypeFor[string](), target.ANY),
		maker.MustMakeDefWithDoc("codemark:registry:doc", reflect.TypeFor[string](), api.OptionDoc{Doc: "some doc"}, target.ANY),
	}
}

func TestInMemory(t *testing.T) {
	reg := InMemory()
	tester := newRegTester(reg)
	defs := defs()
	tests := make([]registryTestCase, 0, len(defs))
	for _, def := range defs {
		tests = append(tests, tester.newTest(def, true))
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			tester.run(t, tc)
		})
	}
}
