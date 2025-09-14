package codemark

import (
	"fmt"
	"testing"
)

func TestOptDocResourcer(t *testing.T) {
	art, err := gen("testdata/proj.go", "")
	if err != nil {
		t.Errorf("unexpected error: %s\n", err)
	}
	fmt.Println(art)
}
