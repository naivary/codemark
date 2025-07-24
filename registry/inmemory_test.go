package registry

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/naivary/codemark/api/doc"
	optionapi "github.com/naivary/codemark/api/option"

	"github.com/naivary/codemark/option"
)

type registryTestCase struct {
	// Name of the test case
	Name string
	// The definition being tested
	Opt *optionapi.Option
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

func (r *registryTester) newTest(opt *optionapi.Option, isValid bool) registryTestCase {
	return registryTestCase{
		Name:    fmt.Sprintf("%s[%s]", opt.Ident, opt.Output),
		IsValid: isValid,
		Opt:     opt,
	}
}

func (r *registryTester) run(t *testing.T, tc registryTestCase) {
	err := r.reg.Define(tc.Opt)
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
	if opt != tc.Opt {
		t.Errorf("definitions are not equal after retrieval. got: %v\n want: %v\n", opt, tc.Opt)
	}
	r.validateDoc(t, tc.Opt, opt)
	t.Logf("test case sucessfull: %s\n", tc.Name)
}

func (r *registryTester) validateDoc(t *testing.T, got, want *optionapi.Option) {
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

func opts() []*optionapi.Option {
	return []*optionapi.Option{
		option.MustMake("codemark:registry:plain", reflect.TypeFor[string](), optionapi.TargetAny),
		option.MustMakeWithDoc(
			"codemark:registry:doc",
			reflect.TypeFor[string](),
			doc.Option{Desc: "some doc"},
			optionapi.TargetAny,
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
