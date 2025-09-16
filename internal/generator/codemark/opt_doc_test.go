package codemark

import (
	"testing"
)

func TestOptDocResourcer(t *testing.T) {
	art, err := gen("testdata/opts.go", "")
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
	_ = art
}
