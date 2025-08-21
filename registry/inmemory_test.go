package registry

import (
	"fmt"
	"reflect"
	"testing"

	docv1 "github.com/naivary/codemark/api/doc/v1"
	optionv1 "github.com/naivary/codemark/api/option/v1"
	regv1 "github.com/naivary/codemark/api/registry/v1"
	"github.com/naivary/codemark/optionutil"
)

type registryTestCase struct {
	// Name of the test case
	Name string
	// The definition being tested
	Opt optionv1.Option
	// Whether the test case is checking correctness or not.
	IsValid bool
}

type registryTester struct {
	reg regv1.Registry
}

func newRegTester(reg regv1.Registry) *registryTester {
	return &registryTester{
		reg: reg,
	}
}

func (r *registryTester) newTest(opt optionv1.Option, isValid bool) registryTestCase {
	return registryTestCase{
		Name:    fmt.Sprintf("%s[%s]", opt.Ident, opt.Type),
		IsValid: isValid,
		Opt:     opt,
	}
}

func (r *registryTester) run(t *testing.T, tc registryTestCase) {
	err := r.reg.Define(&tc.Opt)
	if err == nil && !tc.IsValid {
		t.Errorf("expected an error but err was nil: %s\n", tc.Opt.Ident)
	}
	if err != nil && !tc.IsValid {
		t.Skipf("expected error to occur because it's an invalid test case: %s. Skipping...", err)
	}
	if err != nil && tc.IsValid {
		t.Errorf("could not define definition: %v\n", tc.Opt)
	}
	opt, err := r.reg.Get(tc.Opt.Ident)
	if err != nil {
		t.Errorf("get failed with an error: %s\n", err)
	}
	if opt != &tc.Opt {
		t.Errorf("definitions are not equal after retrieval. got: %v\n want: %v\n", opt, tc.Opt)
	}
	r.validateDoc(t, &tc.Opt, opt)
	t.Logf("test case sucessfull: %s\n", tc.Name)
}

func (r *registryTester) validateDoc(t *testing.T, got, want *optionv1.Option) {
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

func opts() []optionv1.Option {
	return []optionv1.Option{
		optionutil.MustMake(
			"codemark:registry:plain",
			reflect.TypeFor[string](),
			nil, false,
			optionv1.TargetAny),
		optionutil.MustMake(
			"codemark:registry:doc",
			reflect.TypeFor[string](),
			&docv1.Option{Desc: "some doc"},
			false,
			optionv1.TargetAny,
		),
	}
}

func TestInMemory(t *testing.T) {
	reg := InMemory()
	tester := newRegTester(reg)
	opts := opts()
	tests := make([]registryTestCase, 0, len(opts))
	for _, opt := range opts {
		tests = append(tests, tester.newTest(opt, true))
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			tester.run(t, tc)
		})
	}
}
