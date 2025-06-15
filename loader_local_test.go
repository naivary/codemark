package codemark

import (
	"testing"

	sdktesting "github.com/naivary/codemark/sdk/testing"
	"golang.org/x/tools/go/packages"
)

func TestLoaderLocal(t *testing.T) {
	tester, err := sdktesting.NewLoaderTester()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	afs, err := tester.NewFS()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	overlayer := sdktesting.NewInMemOverlayer(afs)
	overlay, err := overlayer.Overlay()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	cfg := &packages.Config{
		Overlay: overlay,
	}
	reg, err := sdktesting.NewDefsSet(NewInMemoryRegistry(), &DefinitionMarker{})
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	mngr, err := NewConvMngr(reg)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	l := NewLocalLoader(mngr, cfg)
	projs, err := l.Load("file=/tmp/funcs.go")
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	for _, s := range projs[0].Funcs {
		t.Logf("name: %s; defs: %v\n", s.Decl.Name, s.Defs)
	}
}
