package testing

import (
	"fmt"
	"testing"

	"github.com/naivary/codemark/definition"
	"github.com/naivary/codemark/sdk"
)

type RegistryTestCase struct {
	// Name of the test case
	Name string
	// The definition being tested
	Def *definition.Definition
	// Whether the test case is checking correctness or not.
	IsValid bool
}

type RegistryTester interface {
	// NewTest returns a new test case
	NewTest(def *definition.Definition, isValid bool) RegistryTestCase

	// ValidTests creates valid test cases for all provided definitions. The test
	// cases are valid in regard of `IsValid` and not the registry interactions.
	ValidTests(defs ...*definition.Definition) []RegistryTestCase

	// Run runs the given test cases and makes all the needed assertions
	Run(t *testing.T, tc RegistryTestCase)
}

var _ RegistryTester = (*registryTester)(nil)

type registryTester struct {
	reg sdk.Registry
}

func NewRegTester(reg sdk.Registry) RegistryTester {
	return &registryTester{
		reg: reg,
	}
}

func (r *registryTester) NewTest(def *definition.Definition, isValid bool) RegistryTestCase {
	return RegistryTestCase{
		Name:    fmt.Sprintf("%s[%s]", def.Ident, def.Output),
		IsValid: isValid,
		Def:     def,
	}
}

func (r *registryTester) ValidTests(defs ...*definition.Definition) []RegistryTestCase {
	tests := make([]RegistryTestCase, 0, len(defs))
	for _, def := range defs {
		tests = append(tests, r.NewTest(def, true))
	}
	return tests
}

func (r *registryTester) Run(t *testing.T, tc RegistryTestCase) {
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

func (r *registryTester) validateDoc(t *testing.T, got, want *definition.Definition) {
	if got.Doc == "" {
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
