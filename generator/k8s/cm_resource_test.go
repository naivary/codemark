package k8s

import (
	"testing"

	"github.com/naivary/codemark"
)

func TestConfigMapResourec(t *testing.T) {
	gen, err := NewGenerator()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	infos, err := codemark.Load(gen.Registry(), "./tests")
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	if err := gen.Generate(infos); err != nil {
		t.Errorf("err occured: %s\n", err)
	}
}
