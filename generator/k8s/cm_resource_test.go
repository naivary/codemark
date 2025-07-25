package k8s

import (
	"testing"

	"github.com/naivary/codemark/loader"
)

func TestConfigMapResource(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	infos, err := loader.Load(gen.Registry(), "./tests")
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	artifacts, err := gen.Generate(infos)
	if err != nil {
		t.Errorf("err occured: %s\n", err)
	}
	_ = artifacts
}
