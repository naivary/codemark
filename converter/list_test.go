package converter

import (
	"testing"

	"github.com/naivary/codemark/converter/convertertest"
	"github.com/naivary/codemark/registry/registrytest"
)

func newManager() *Manager {
	reg, err := registrytest.NewRegistry(registrytest.NewOptsSet())
	if err != nil {
		panic(err)
	}
	mngr, err := NewManager(reg)
	if err != nil {
		panic(err)
	}
	return mngr
}

func TestListConverter(t *testing.T) {
	conv := NewList(newManager())
	tester, err := convertertest.NewTester(conv)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	cases, err := validCasesFor(conv)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tester.Run(t, tc)
		})
	}
}
