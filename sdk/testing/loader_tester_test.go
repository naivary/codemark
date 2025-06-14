package testing

import "testing"

func TestRandField(t *testing.T) {
	s := randStruct()
	t.Log(s.name)
	for _, m := range s.markers {
		t.Log(m)
	}
	for _, f := range s.fields {
		t.Log(f.f.Name)
	}
}
