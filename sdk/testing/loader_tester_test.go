package testing

import (
	"io"
	"strings"
	"testing"
)

func TestLoaderTester(t *testing.T) {
	tester, err := NewLoaderTester()
	if err != nil {
		t.Errorf("err occured: %v\n", err)
	}
	r, err := tester.NewFile()
	if err != nil {
		t.Errorf("err occured: %v\n", err)
	}
	var b strings.Builder
	_, err = io.Copy(&b, r)
	if err != nil {
		t.Errorf("err occured: %v\n", err)
	}
	t.Log(b.String())
}
